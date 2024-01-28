package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"

	"github.com/wetor/AnimeGo/internal/web"
	webapi "github.com/wetor/AnimeGo/internal/web/api"
	_ "github.com/wetor/AnimeGo/internal/web/docs"
	"github.com/wetor/AnimeGo/internal/web/websocket"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/cmd/common"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/animego/anisource/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/internal/pkg/torrent"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/wire"
)

const (
	DefaultConfigFile = "data/animego.yaml"
)

var (
	ctx, cancel  = context.WithCancel(context.Background())
	configFile   string
	debug        bool
	webapiEnable bool
	backupConfig bool
)

func main() {
	common.PrintInfo()
	flag.StringVar(&configFile, "config", DefaultConfigFile, "配置文件路径；配置文件中的相对路径均是相对与程序的位置")
	flag.BoolVar(&debug, "debug", false, "Debug模式，将会显示更多的日志")
	flag.BoolVar(&webapiEnable, "web", true, "启用Web API")
	flag.BoolVar(&backupConfig, "backup", true, "配置文件升级前是否备份")
	flag.Parse()

	tmp := os.Getenv("ANIMEGO_CONFIG")
	if len(tmp) > 0 {
		configFile = tmp
	}
	tmp = os.Getenv("ANIMEGO_DEBUG")
	if len(tmp) > 0 {
		debug = utils.String2Bool(tmp)
	}
	tmp = os.Getenv("ANIMEGO_WEB")
	if len(tmp) > 0 {
		webapiEnable = utils.String2Bool(tmp)
	}
	tmp = os.Getenv("ANIMEGO_CONFIG_BACKUP")
	if len(tmp) > 0 {
		backupConfig = utils.String2Bool(tmp)
	}

	common.RegisterExit(doExit)
	Main()
}

func doExit() {
	log.Infof("正在退出...")
	cancel()
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}

func Main() {
	var err error
	var wg sync.WaitGroup
	var bangumiCacheMutex sync.Mutex

	configFile = xpath.Abs(configFile)
	// 初始化默认配置、升级配置
	if utils.IsExist(configFile) {
		configs.InitUpdateConfig(configFile, backupConfig)
	} else {
		configs.InitDefaultConfig(configFile)
	}
	// 解析环境变量文本并写入配置文件
	configs.InitEnvConfig(configFile, configFile)

	// ===============================================================================================================
	// 载入配置文件
	config := configs.Load(configFile)
	constant.Init(&constant.Options{
		DataPath: xpath.P(config.DataPath),
	})
	// 创建子文件夹
	config.InitDir()
	// 检查参数限制
	config.Check()
	// 释放资源
	assets.WritePlugins(assets.Dir, path.Join(xpath.P(config.DataPath), assets.Dir), true)

	// ===============================================================================================================
	// 初始化日志
	out, notify := logger.NewLogNotify()
	logger.Init(&logger.Options{
		File:    constant.LogFile,
		Debug:   debug,
		Context: ctx,
		WG:      &wg,
		Out:     out,
	})

	// 初始化request
	request.Init(&request.Options{
		UserAgent: fmt.Sprintf("%s/AnimeGo (%s)", os.Getenv("ANIMEGO_VERSION"), constant.AnimeGoGithub),
		Proxy:     config.Proxy(),
		Timeout:   config.Advanced.Request.TimeoutSecond,
		Retry:     config.Advanced.Request.RetryNum,
		RetryWait: config.Advanced.Request.RetryWaitSecond,
		Debug:     debug,
		Host: map[string]*request.HostOptions{
			constant.BangumiHost: {
				Redirect: config.Advanced.Source.Bangumi.Redirect,
			},
			constant.MikanHost: {
				Redirect: config.Advanced.Source.Mikan.Redirect,
				Header: map[string]string{
					constant.MikanAuthCookie: config.Advanced.Source.Mikan.Cookie,
				},
			},
			constant.ThemoviedbHost: {
				Redirect: config.Advanced.Source.Themoviedb.Redirect,
				Params: map[string]string{
					constant.ThemoviedbApiKey: config.Advanced.Source.Themoviedb.ApiKey,
				},
			},
		},
	})
	// 初始化torrent
	torrent.Init(&torrent.Options{
		TempPath: constant.TempPath,
	})
	// ===============================================================================================================
	// 载入AnimeGo数据库（缓存）
	bolt := cache.NewBolt()
	bolt.Open(constant.CacheFile)

	// 载入Bangumi Archive数据库
	bangumiCache := cache.NewBolt()
	bangumiCache.Open(constant.BangumiCacheFile)

	bgmOpts := &bangumi.Options{
		Cache:            bolt,
		CacheTime:        int64(config.Advanced.Cache.BangumiCacheHour * 60 * 60),
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &bangumiCacheMutex,
	}
	mikanOpts := &mikan.Options{
		Cache:     bolt,
		CacheTime: int64(config.Advanced.Cache.MikanCacheHour * 60 * 60),
	}
	tmdbOpts := &themoviedb.Options{
		Cache:     bolt,
		CacheTime: int64(config.Advanced.Cache.ThemoviedbCacheHour * 60 * 60),
	}
	// ===============================================================================================================
	// 初始化插件 gpython
	plugin.Init(&plugin.Options{
		Path:  constant.PluginPath,
		Debug: debug,
		Feed:  feed.NewRss(),
		Mikan: wire.GetMikanData(mikanOpts),
	})

	// ===============================================================================================================
	// 初始化并连接下载器
	clientOpts := &models.ClientOptions{
		WG:  &wg,
		Ctx: ctx,

		Url:      config.Setting.Client.Url,
		Username: config.Setting.Client.Username,
		Password: config.Setting.Client.Password,

		DownloadPath:         config.Setting.Client.DownloadPath,
		SeedingTimeMinute:    config.Advanced.Client.SeedingTimeMinute,
		ConnectTimeoutSecond: config.Advanced.Client.ConnectTimeoutSecond,
		CheckTimeSecond:      config.Advanced.Client.CheckTimeSecond,
		RetryConnectNum:      config.Advanced.Client.RetryConnectNum,
	}
	clientSrv := wire.GetClient(config.Setting.Client.Client, clientOpts, bolt)
	clientSrv.Start()

	// ===============================================================================================================
	// 第一个启用的rename插件
	var renamePlugin *models.Plugin
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Rename) {
		if p.Enable {
			renamePlugin = &p
			break
		}
	}
	// 初始化rename
	renameSrv := wire.GetRenamer(&models.RenamerOptions{
		WG:            &wg,
		RefreshSecond: config.RefreshSecond,
	}, renamePlugin)
	// 启动rename
	renameSrv.Start(ctx)

	// ===============================================================================================================
	// 初始化database配置
	databaseInst, err := wire.GetDatabase(&models.DatabaseOptions{
		SavePath: xpath.P(config.SavePath),
	}, bolt)
	if err != nil {
		panic(err)
	}

	// ===============================================================================================================
	// 初始化downloader

	downloadCallback := &models.Callback{}
	downloaderSrv := wire.GetDownloader(&models.DownloaderOptions{
		RefreshSecond:          config.RefreshSecond,
		Category:               config.Category,
		Tag:                    config.Tag,
		AllowDuplicateDownload: config.Download.AllowDuplicateDownload,
		WG:                     &wg,
	}, clientSrv, &models.NotifierOptions{
		DownloadPath: xpath.P(config.DownloadPath),
		SavePath:     xpath.P(config.SavePath),
		Rename:       config.Advanced.Download.Rename,
		Callback:     downloadCallback,
	}, databaseInst, renameSrv)

	downloadCallback.Func = func(data any) error {
		return downloaderSrv.Delete(data.(string))
	}
	// 启动downloader
	downloaderSrv.Start(ctx)

	// ===============================================================================================================
	// 第一个启用的parser插件
	var parsePlugin *models.Plugin
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Parser) {
		if p.Enable {
			parsePlugin = &p
			break
		}
	}
	// 初始化filter

	filterSrv := wire.GetFilter(&models.FilterOptions{
		DelaySecond: config.Advanced.Feed.DelaySecond,
	}, downloaderSrv, &models.ParserOptions{
		TMDBFailSkip:           config.Default.TMDBFailSkip,
		TMDBFailUseTitleSeason: config.Default.TMDBFailUseTitleSeason,
		TMDBFailUseFirstSeason: config.Default.TMDBFailUseFirstSeason,
	}, parsePlugin, mikanOpts, bgmOpts, tmdbOpts)
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Filter) {
		filterSrv.Add(&p)
	}
	// ===============================================================================================================
	// 初始化定时任务
	scheduleSrv := schedule.NewSchedule(&models.ScheduleOptions{
		WG: &wg,
	})
	// 添加定时任务
	err = scheduleSrv.Add(&schedule.AddTaskOptions{
		Name:     "bangumi",
		StartRun: true,
		Task: schedule.NewBangumiTask(&schedule.BangumiOptions{
			Cache:      bangumiCache,
			CacheMutex: &bangumiCacheMutex,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = scheduleSrv.Add(&schedule.AddTaskOptions{
		Name:     "database",
		StartRun: false,
		Task: schedule.NewRefreshTask(&schedule.RefreshOptions{
			Database: databaseInst,
			Cron:     config.Advanced.Database.RefreshDatabaseCron,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = schedule.AddScheduleTasks(scheduleSrv, configs.ConvertPluginInfo(config.Plugin.Schedule))
	if err != nil {
		panic(err)
	}
	err = feed.AddFeedTasks(scheduleSrv, configs.ConvertPluginInfo(config.Plugin.Feed), filterSrv, ctx)
	if err != nil {
		panic(err)
	}
	// 启动化定时任务
	scheduleSrv.Start(ctx)

	// ===============================================================================================================
	if webapiEnable {
		// 初始化Web API
		web.Init(&web.Options{
			ApiOptions: &webapi.Options{
				Ctx:                  ctx,
				AccessKey:            config.WebApi.AccessKey,
				Cache:                bolt,
				Config:               config,
				BangumiCache:         bangumiCache,
				BangumiCacheLock:     &bangumiCacheMutex,
				FilterManager:        filterSrv,
				DatabaseCacheDeleter: databaseInst,
			},
			WebSocketOptions: &websocket.Options{
				WG:     &wg,
				Notify: notify,
			},
			Host:  config.WebApi.Host,
			Port:  config.WebApi.Port,
			WG:    &wg,
			Debug: debug,
		})
		// 启动Web API
		web.Run(ctx)
	}

	// 等待程序运行结束
	wg.Wait()
}

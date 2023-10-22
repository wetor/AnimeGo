package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/cmd/common"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	anidataBangumi "github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	anidataMikan "github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	anidataThemoviedb "github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/animego/downloader"
	feedPlugin "github.com/wetor/AnimeGo/internal/animego/feed/plugin"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/animego/parser"
	parserPlugin "github.com/wetor/AnimeGo/internal/animego/parser/plugin"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/internal/web"
	webapi "github.com/wetor/AnimeGo/internal/web/api"
	_ "github.com/wetor/AnimeGo/internal/web/docs"
	"github.com/wetor/AnimeGo/internal/web/websocket"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/client/transmission"
	pkgLog "github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/torrent"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
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

	WG                sync.WaitGroup
	BangumiCacheMutex sync.Mutex
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
	pkgLog.Infof("正在退出...")
	cancel()
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}

func Main() {
	var err error
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
		WG:      &WG,
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
	})
	// 初始化torrent
	torrent.Init(&torrent.Options{
		TempPath: constant.TempPath,
	})
	// ===============================================================================================================
	// 初始化插件 gpython
	plugin.Init(&plugin.Options{
		Path:  constant.PluginPath,
		Debug: debug,
	})
	// 载入AnimeGo数据库（缓存）
	bolt := cache.NewBolt()
	bolt.Open(constant.CacheFile)

	// 载入Bangumi Archive数据库
	bangumiCache := cache.NewBolt()
	bangumiCache.Open(constant.BangumiCacheFile)

	// ===============================================================================================================
	// 初始化并连接下载器
	//qbittorrentSrv := qbittorrent.NewQBittorrent(&qbittorrent.Options{
	//	Url:                  config.Setting.Client.QBittorrent.Url,
	//	Username:             config.Setting.Client.QBittorrent.Username,
	//	Password:             config.Setting.Client.QBittorrent.Password,
	//	DownloadPath:         config.Setting.Client.QBittorrent.DownloadPath,
	//	ConnectTimeoutSecond: config.Advanced.Client.ConnectTimeoutSecond,
	//	CheckTimeSecond:      config.Advanced.Client.CheckTimeSecond,
	//	RetryConnectNum:      config.Advanced.Client.RetryConnectNum,
	//	WG:                   &WG,
	//	Ctx:                  ctx,
	//})
	//qbittorrentSrv.Start()

	transmissionSrv := transmission.NewTransmission(&transmission.Options{
		Url:                  config.Setting.Client.QBittorrent.Url,
		Username:             config.Setting.Client.QBittorrent.Username,
		Password:             config.Setting.Client.QBittorrent.Password,
		DownloadPath:         config.Setting.Client.QBittorrent.DownloadPath,
		ConnectTimeoutSecond: config.Advanced.Client.ConnectTimeoutSecond,
		CheckTimeSecond:      config.Advanced.Client.CheckTimeSecond,
		RetryConnectNum:      config.Advanced.Client.RetryConnectNum,
		WG:                   &WG,
		Ctx:                  ctx,
	})
	transmissionSrv.Start()
	// ===============================================================================================================
	// 初始化anisource配置
	anisource.Init(&anisource.Options{
		AniDataOptions: &anidata.Options{
			Cache: bolt,
			CacheTime: map[string]int64{
				anidataMikan.Bucket:      int64(config.Advanced.Cache.MikanCacheHour * 60 * 60),
				anidataBangumi.Bucket:    int64(config.Advanced.Cache.BangumiCacheHour * 60 * 60),
				anidataThemoviedb.Bucket: int64(config.Advanced.Cache.ThemoviedbCacheHour * 60 * 60),
			},
			BangumiCache:       bangumiCache,
			BangumiCacheLock:   &BangumiCacheMutex,
			RedirectMikan:      config.Advanced.Redirect.Mikan,
			RedirectBangumi:    config.Advanced.Redirect.Bangumi,
			RedirectThemoviedb: config.Advanced.Redirect.Themoviedb,
		},
	})

	// ===============================================================================================================
	// 初始化renamer配置
	renamer.Init(&renamer.Options{
		WG:            &WG,
		RefreshSecond: config.RefreshSecond,
	})
	// 第一个启用的rename插件
	var rename api.RenamerPlugin
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Rename) {
		if p.Enable {
			rename = renamerPlugin.NewRenamePlugin(&p)
			break
		}
	}
	// 初始化rename
	renameSrv := renamer.NewManager(rename)
	// 启动rename
	renameSrv.Start(ctx)

	// ===============================================================================================================
	// 初始化database配置
	database.Init(&database.Options{
		DownloaderConf: database.DownloaderConf{
			RefreshSecond:          config.RefreshSecond,
			DownloadPath:           xpath.P(config.DownloadPath),
			SavePath:               xpath.P(config.SavePath),
			Category:               config.Category,
			Tag:                    config.Tag,
			AllowDuplicateDownload: config.Download.AllowDuplicateDownload,
			SeedingTimeMinute:      config.Download.SeedingTimeMinute,
			Rename:                 config.Advanced.Download.Rename,
		},
	})
	downloadCallback := &database.Callback{}
	databaseSrv, err := database.NewDatabase(bolt, renameSrv, downloadCallback)
	if err != nil {
		panic(err)
	}

	// ===============================================================================================================
	// 初始化downloader配置
	downloader.Init(&downloader.Options{
		RefreshSecond:          config.RefreshSecond,
		Category:               config.Category,
		Tag:                    config.Tag,
		AllowDuplicateDownload: config.Download.AllowDuplicateDownload,
		SeedingTimeMinute:      config.Download.SeedingTimeMinute,
		WG:                     &WG,
	})
	// 初始化downloader
	downloaderSrv := downloader.NewManager(transmissionSrv, databaseSrv, databaseSrv)
	downloadCallback.Renamed = func(data any) error {
		return downloaderSrv.Delete(data.(string))
	}
	// 启动downloader
	downloaderSrv.Start(ctx)

	// ===============================================================================================================
	// 初始化parser配置
	parser.Init(&parser.Options{
		TMDBFailSkip:           config.Default.TMDBFailSkip,
		TMDBFailUseTitleSeason: config.Default.TMDBFailUseTitleSeason,
		TMDBFailUseFirstSeason: config.Default.TMDBFailUseFirstSeason,
	})
	// 第一个启用的parser插件
	var parse api.ParserPlugin
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Parser) {
		if p.Enable {
			parse = parserPlugin.NewParserPlugin(&p, false)
			break
		}
	}
	// 初始化parser
	bgmSource := bangumi.NewBangumiSource(config.Setting.Key.Themoviedb)
	mikanSource := mikan.NewMikanSource(bgmSource)
	parserSrv := parser.NewManager(parse, mikanSource, bgmSource)

	// ===============================================================================================================
	// 初始化filter配置
	filter.Init(&filter.Options{
		DelaySecond: config.Advanced.Feed.DelaySecond,
	})
	// 初始化filter
	filterSrv := filter.NewManager(downloaderSrv, parserSrv)
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Filter) {
		filterSrv.Add(&p)
	}
	// ===============================================================================================================
	// 初始化定时任务
	scheduleSrv := schedule.NewSchedule(&schedule.Options{
		WG: &WG,
	})
	// 添加定时任务
	err = scheduleSrv.Add(&schedule.AddTaskOptions{
		Name:     "bangumi",
		StartRun: true,
		Task: task.NewBangumiTask(&task.BangumiOptions{
			Cache:      bangumiCache,
			CacheMutex: &BangumiCacheMutex,
		}),
	})
	if err != nil {
		panic(err)
	}
	err = schedule.AddScheduleTasks(scheduleSrv, configs.ConvertPluginInfo(config.Plugin.Schedule))
	if err != nil {
		panic(err)
	}
	err = feedPlugin.AddFeedTasks(scheduleSrv, configs.ConvertPluginInfo(config.Plugin.Feed), filterSrv, ctx)
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
				BangumiCacheLock:     &BangumiCacheMutex,
				FilterManager:        filterSrv,
				DatabaseCacheDeleter: databaseSrv,
			},
			WebSocketOptions: &websocket.Options{
				WG:     &WG,
				Notify: notify,
			},
			Host:  config.WebApi.Host,
			Port:  config.WebApi.Port,
			WG:    &WG,
			Debug: debug,
		})
		// 启动Web API
		web.Run(ctx)
	}

	// 等待程序运行结束
	WG.Wait()
}

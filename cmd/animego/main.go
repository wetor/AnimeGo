package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/cmd/common"
	"github.com/wetor/AnimeGo/configs"
	_ "github.com/wetor/AnimeGo/docs"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	anidataBangumi "github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	anidataMikan "github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	anidataThemoviedb "github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	feedPlugin "github.com/wetor/AnimeGo/internal/animego/feed/plugin"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/animego/manager"
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
	"github.com/wetor/AnimeGo/internal/web/websocket"
	"github.com/wetor/AnimeGo/pkg/cache"
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

	WG                sync.WaitGroup
	BangumiCacheMutex sync.Mutex
)

func main() {
	common.PrintInfo()

	flag.StringVar(&configFile, "config", DefaultConfigFile, "配置文件路径；配置文件中的相对路径均是相对与程序的位置")
	flag.BoolVar(&debug, "debug", false, "Debug模式，将会显示更多的日志")
	flag.BoolVar(&webapiEnable, "web", true, "启用Web API，默认启用")
	flag.Parse()

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

func InitDefaultConfig() {
	if utils.IsExist(configFile) {
		// 尝试升级配置文件
		if configs.UpdateConfig(configFile, true) {
			os.Exit(0)
		}
		return
	}
	log.Printf("未找到配置文件（%s），开始初始化默认配置\n", configFile)
	conf := configs.DefaultConfig()
	if utils.IsExist(conf.Setting.DataPath) {
		log.Printf("默认data_path文件夹（%s）已存在，无法完成初始化\n", conf.Setting.DataPath)
		os.Exit(0)
	}
	err := utils.CreateMutiDir(conf.Setting.DataPath)
	if err != nil {
		panic(err)
	}
	err = configs.DefaultFile(DefaultConfigFile)
	if err != nil {
		panic(err)
	}

	InitDefaultAssets(conf.DataPath, true)

	log.Printf("初始化默认配置完成（%s）\n", conf.Setting.DataPath)
	log.Println("请设置配置后重新启动")
	os.Exit(0)
}

func InitDefaultAssets(dataPath string, skip bool) {
	assets.WritePlugins(assets.Dir, xpath.Join(dataPath, assets.Dir), skip)
}

func Main() {
	configFile = xpath.Abs(configFile)
	// 初始化默认配置、升级配置
	InitDefaultConfig()

	// ===============================================================================================================
	// 载入配置文件
	config := configs.Init(configFile)
	// 检查参数限制
	config.Check()
	constant.Init(&constant.Options{
		DataPath: config.DataPath,
	})
	config.InitDir()

	// 释放资源
	InitDefaultAssets(config.DataPath, true)

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
	downloader.Init(&downloader.Options{
		ConnectTimeoutSecond: config.Advanced.Client.ConnectTimeoutSecond,
		CheckTimeSecond:      config.Advanced.Client.CheckTimeSecond,
		RetryConnectNum:      config.Advanced.Client.RetryConnectNum,
		WG:                   &WG,
	})
	qbtConf := config.Setting.Client.QBittorrent
	qbittorrentSrv := qbittorrent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password)
	qbittorrentSrv.Start(ctx)

	// ===============================================================================================================
	// 初始化anisource配置
	anisource.Init(&anisource.Options{
		Options: &anidata.Options{
			Cache: bolt,
			CacheTime: map[string]int64{
				anidataMikan.Bucket:      int64(config.Advanced.Cache.MikanCacheHour * 60 * 60),
				anidataBangumi.Bucket:    int64(config.Advanced.Cache.BangumiCacheHour * 60 * 60),
				anidataThemoviedb.Bucket: int64(config.Advanced.Cache.ThemoviedbCacheHour * 60 * 60),
			},
			BangumiCache:     bangumiCache,
			BangumiCacheLock: &BangumiCacheMutex,
		},
	})

	// ===============================================================================================================
	// 初始化renamer配置
	renamer.Init(&renamer.Options{
		WG:                &WG,
		UpdateDelaySecond: config.UpdateDelaySecond,
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
	// 初始化manager配置
	manager.Init(&manager.Options{
		Downloader: manager.Downloader{
			UpdateDelaySecond:      config.UpdateDelaySecond,
			DownloadPath:           config.DownloadPath,
			SavePath:               config.SavePath,
			Category:               config.Category,
			Tag:                    config.Tag,
			AllowDuplicateDownload: config.Download.AllowDuplicateDownload,
			SeedingTimeMinute:      config.Download.SeedingTimeMinute,
			Rename:                 config.Advanced.Download.Rename,
		},
		WG: &WG,
	})
	// 初始化manager
	managerSrv := manager.NewManager(qbittorrentSrv, bolt, renameSrv)
	// 启动manager
	managerSrv.Start(ctx)
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
	parserSrv := parser.NewManager(parse, &mikan.Mikan{ThemoviedbKey: config.Setting.Key.Themoviedb})

	// ===============================================================================================================
	// 初始化filter配置
	filter.Init(&filter.Options{
		DelaySecond: config.Advanced.Feed.DelaySecond,
	})
	// 初始化filter
	filterSrv := filter.NewManager(managerSrv, parserSrv)
	for _, p := range configs.ConvertPluginInfo(config.Plugin.Filter) {
		filterSrv.Add(&p)
	}
	// ===============================================================================================================
	// 初始化定时任务
	scheduleSrv := schedule.NewSchedule(&schedule.Options{
		WG: &WG,
	})
	// 添加定时任务
	scheduleSrv.Add(&schedule.AddTaskOptions{
		Name:     "bangumi",
		StartRun: true,
		Task: task.NewBangumiTask(&task.BangumiOptions{
			Cache:      bangumiCache,
			CacheMutex: &BangumiCacheMutex,
		}),
	})
	schedule.AddScheduleTasks(scheduleSrv, configs.ConvertPluginInfo(config.Plugin.Schedule))
	feedPlugin.AddFeedTasks(scheduleSrv, configs.ConvertPluginInfo(config.Plugin.Feed), filterSrv, ctx)
	// 启动化定时任务
	scheduleSrv.Start(ctx)

	// ===============================================================================================================
	if webapiEnable {
		// 初始化Web API
		web.Init(&web.Options{
			ApiOptions: &webapi.Options{
				Ctx:                           ctx,
				AccessKey:                     config.WebApi.AccessKey,
				Cache:                         bolt,
				Config:                        config,
				BangumiCache:                  bangumiCache,
				BangumiCacheLock:              &BangumiCacheMutex,
				FilterManager:                 filterSrv,
				DownloaderManagerCacheDeleter: managerSrv,
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

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/configs"
	_ "github.com/wetor/AnimeGo/docs"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/internal/process/animego"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/internal/web"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/third_party/gpython"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

const (
	AnimeGoVersion       = "0.5.2"
	AnimeGoConfigVersion = "1.1.0"
	AnimeGoGithub        = "https://github.com/wetor/AnimeGo"

	DefaultConfigFile = "./data/animego.yaml"
)

var ctx, cancel = context.WithCancel(context.Background())
var configFile string
var debug bool

func init() {
	var err error
	err = os.Setenv("ANIMEGO_VERSION", AnimeGoVersion)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("ANIMEGO_CONFIG_VERSION", AnimeGoConfigVersion)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("ANIMEGO_GITHUB", AnimeGoGithub)
	if err != nil {
		panic(err)
	}
}

func main() {
	printInfo()

	flag.StringVar(&configFile, "config", DefaultConfigFile, "配置文件路径；配置文件中的相对路径均是相对与程序的位置")
	flag.BoolVar(&debug, "debug", false, "Debug模式，将会显示更多的日志")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				zap.S().Infof("收到退出信号: %v", s)
				doExit()
			default:
				zap.S().Infof("收到其他信号: %v", s)
			}
		}
	}()
	Main(ctx)
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

	InitDefaultAssets(conf)

	log.Printf("初始化默认配置完成（%s）\n", conf.Setting.DataPath)
	log.Println("请设置配置后重新启动")
	os.Exit(0)
}

func InitDefaultAssets(conf *configs.Config) {
	utils.CopyDir(assets.Plugin, "plugin", path.Join(conf.Setting.DataPath, "plugin"), true, true)
}

func doExit() {
	zap.S().Infof("正在退出...")
	cancel()
	go func() {
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()
}

func printInfo() {
	fmt.Println(`--------------------------------------------------
    ___            _                   ______     
   /   |   ____   (_)____ ___   ___   / ____/____ 
  / /| |  / __ \ / // __ \__ \ / _ \ / / __ / __ \
 / ___ | / / / // // / / / / //  __// /_/ // /_/ /
/_/  |_|/_/ /_//_//_/ /_/ /_/ \___/ \____/ \____/
    `)
	fmt.Printf("AnimeGo v%s\n", os.Getenv("ANIMEGO_VERSION"))
	fmt.Printf("AnimeGo config v%s\n", os.Getenv("ANIMEGO_CONFIG_VERSION"))
	fmt.Printf("%s\n", os.Getenv("ANIMEGO_GITHUB"))
	fmt.Println("--------------------------------------------------")
}

func Main(ctx context.Context) {
	InitDefaultConfig()

	config := configs.Init(configFile)
	config.InitDir()

	InitDefaultAssets(config)

	logger.Init(&logger.InitOptions{
		File:    config.Advanced.Path.LogFile,
		Debug:   debug,
		Context: ctx,
	})

	bolt := cache.NewBolt()
	bolt.Open(config.Advanced.Path.DbFile)
	bangumiCache := cache.NewBolt()
	bangumiCache.Open(path.Join(path.Dir(config.Advanced.Path.DbFile), "bolt_sub.db"))

	store.Init(&store.InitOptions{
		Config:       config,
		Cache:        bolt,
		BangumiCache: bangumiCache,
	})

	request.Init(&request.InitOptions{
		Proxy:     store.Config.Proxy(),
		Timeout:   store.Config.Advanced.Request.TimeoutSecond,
		Retry:     store.Config.Advanced.Request.RetryNum,
		RetryWait: store.Config.Advanced.Request.RetryWaitSecond,
		Debug:     debug,
	})
	gpython.Init()

	public.Init(&public.Options{
		PluginPath: path.Join(store.Config.DataPath, "plugin"),
	})

	store.Process = animego.NewAnimeGo()
	store.Process.Run(ctx)

	web.Init(&web.InitOptions{
		Debug: debug,
	})

	web.Run(ctx)
	store.WG.Wait()
}

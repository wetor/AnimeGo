package main

import (
	"AnimeGo/configs"
	"AnimeGo/internal/cache"
	"AnimeGo/internal/logger"
	"AnimeGo/internal/process/mikan"
	"AnimeGo/internal/store"
	"AnimeGo/internal/web"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

var ctx, cancel = context.WithCancel(context.Background())
var configFile string
var debug bool

func main() {
	flag.StringVar(&configFile, "config", "data/config/animego.yaml", "配置文件路径；配置文件中的相对路径均是相对与程序的位置")
	flag.BoolVar(&debug, "debug", true, "Debug模式，将会输出更多的日志")

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

func doExit() {
	zap.S().Infof("正在退出...")
	cancel()
}

func Main(ctx context.Context) {

	config := configs.Init(configFile)
	config.InitDir()

	logger.Init(&logger.InitOptions{
		File:    config.LogFile,
		Debug:   debug,
		Context: ctx,
	})

	bolt := cache.NewBolt()
	bolt.Open(config.Setting.DbFile)

	store.Init(&store.InitOptions{
		Config: config,
		Cache:  bolt,
	})

	store.Process = mikan.NewMikan()
	store.Process.Run(ctx)
	web.Run(ctx)
	store.WG.Wait()
}

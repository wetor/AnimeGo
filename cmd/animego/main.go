package main

import (
	"GoBangumi/process"
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"context"
	"flag"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var ctx, cancel = context.WithCancel(context.Background())
var configFile string

func main() {
	flag.StringVar(&configFile, "config", "data/config/conf.yaml", "配置文件路径；配置文件中的相对路径均是相对与程序的位置")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for s := range sigs {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT:
				zap.S().Infof("收到退出信号: %v", s)
				doExit()
			case syscall.SIGUSR1:
				zap.S().Infof("收到USR1信号: %v", s)
			case syscall.SIGUSR2:
				zap.S().Infof("收到USR2信号: %v", s)
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
	_ = zap.S().Sync()
}

func Main(ctx context.Context) {
	logger.Init()
	defer logger.Flush()
	store.Init(&store.InitOptions{
		ConfigFile: configFile,
	})
	store.WG.Add(2)
	m := process.NewMikan()
	m.Run(ctx)
	store.WG.Wait()
}

package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/process/mikan"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/internal/web"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/request"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

const (
	AnimeGoVersion = "0.3.2"
	AnimeGoGithub  = "https://github.com/wetor/AnimeGo"
)

var ctx, cancel = context.WithCancel(context.Background())
var configFile string
var debug bool

var rootPath string
var replace bool

func init() {
	var err error
	err = os.Setenv("animego_version", AnimeGoVersion)
	if err != nil {
		panic(err)
	}
	err = os.Setenv("animego_github", AnimeGoGithub)
	if err != nil {
		panic(err)
	}
}

func main() {
	printInfo()

	flag.StringVar(&configFile, "config", "data/config/animego.yaml", "配置文件路径；配置文件中的相对路径均是相对与程序的位置")
	flag.BoolVar(&debug, "debug", false, "Debug模式，将会输出更多的日志")

	flag.StringVar(&rootPath, "init-path", "", "[初始化]输出资源/配置文件到的根目录")
	flag.BoolVar(&replace, "init-replace", false, "[初始化]输出资源/配置文件时是否自动替换")
	flag.Parse()
	if len(rootPath) > 0 {
		copyDir(assets.Plugin, "plugin", path.Join(rootPath, "plugin"), replace)
		copyDir(assets.Config, "config", path.Join(rootPath, "config"), replace)
		return
	}

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

func printInfo() {
	fmt.Printf("github.com/wetor/AnimeGo %s (%s)\n", os.Getenv("animego_version"), os.Getenv("animego_github"))
}

func Main(ctx context.Context) {

	config := configs.Init(configFile)

	logger.Init(&logger.InitOptions{
		File:    config.Advanced.Path.LogFile,
		Debug:   debug,
		Context: ctx,
	})

	bolt := cache.NewBolt()
	bolt.Open(config.Advanced.Path.DbFile)

	store.Init(&store.InitOptions{
		Config: config,
		Cache:  bolt,
	})

	request.Init(&request.InitOptions{
		Proxy:     store.Config.Proxy(),
		Timeout:   store.Config.Advanced.Request.TimeoutSecond,
		Retry:     store.Config.Advanced.Request.RetryNum,
		RetryWait: store.Config.Advanced.Request.RetryWaitSecond,
		Debug:     debug,
	})

	store.Process = mikan.NewMikan()
	store.Process.Run(ctx)

	web.Init(&web.InitOptions{
		Debug: debug,
	})

	web.Run(ctx)
	store.WG.Wait()
}

func copyDir(fs embed.FS, src, dst string, replace bool) {
	files, err := fs.ReadDir(src)
	if err != nil {
		panic(err)
	}

	err = utils.CreateMutiDir(dst)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		srcPath := path.Join(src, file.Name())
		dstPath := path.Join(dst, file.Name())
		if file.IsDir() {
			copyDir(fs, srcPath, dstPath, replace)
			continue
		}
		fileContent, err := fs.ReadFile(srcPath)
		if err != nil {
			panic(err)
		}
		if !replace && utils.IsExist(dstPath) {
			fmt.Printf("文件[%s]已存在，是否替换[y(yes)/n(no)]: ", dstPath)
			if !scanYesNo() {
				continue
			}
		}
		if err := os.WriteFile(dstPath, fileContent, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func scanYesNo() bool {
	var s string
	_, err := fmt.Scanln(&s)
	if err != nil {
		panic(err)
	}
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

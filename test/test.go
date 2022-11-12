package test

import (
	"context"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/anisource"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/request"
)

func TestInit() {
	config := configs.Init("/Users/wetor/GoProjects/AnimeGo/data/animego.yaml")
	config.InitDir()
	logger.Init(&logger.InitOptions{
		File:    config.Advanced.Path.LogFile,
		Debug:   true,
		Context: context.Background(),
	})

	bolt := cache.NewBolt()
	bolt.Open(config.Advanced.Path.DbFile)
	anisource.Init(&anisource.Options{Cache: bolt})
	store.Init(&store.InitOptions{
		Config: config,
		Cache:  bolt,
	})
	request.Init(&request.InitOptions{
		Proxy:     store.Config.Proxy(),
		Timeout:   store.Config.Advanced.Request.TimeoutSecond,
		Retry:     store.Config.Advanced.Request.RetryNum,
		RetryWait: store.Config.Advanced.Request.RetryWaitSecond,
		Debug:     true,
	})
}

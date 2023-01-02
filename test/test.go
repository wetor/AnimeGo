package test

import (
	"context"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/request"
	"path"
)

func TestInit() {
	ctx := context.Background()
	debug := true

	config := configs.Init("/Users/wetor/GoProjects/AnimeGo/data/animego.yaml")
	config.InitDir()

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
}

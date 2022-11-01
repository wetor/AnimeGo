package test

import (
	"context"
	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/cache"
)

func TestInit() {
	config := configs.Init("/Users/wetor/GoProjects/AnimeGo/data/config/animego.yaml")
	config.InitDir()

	logger.Init(&logger.InitOptions{
		File:    config.Advanced.Path.LogFile,
		Debug:   true,
		Context: context.Background(),
	})

	bolt := cache.NewBolt()
	bolt.Open(config.Advanced.Path.DbFile)

	store.Init(&store.InitOptions{
		Config: config,
		Cache:  bolt,
	})
}

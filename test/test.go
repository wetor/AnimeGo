package test

import (
	"AnimeGo/configs"
	"AnimeGo/internal/cache"
	"AnimeGo/internal/logger"
	"AnimeGo/internal/store"
	"context"
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

package store

import (
	"GoBangumi/configs"
	"GoBangumi/internal/cache"
	"GoBangumi/utils"
	"go.uber.org/zap"
	"sync"
)

var (
	Cache  cache.Cache
	Config *configs.Config
	WG     sync.WaitGroup
)

type InitOptions struct {
	Cache      cache.Cache
	ConfigFile string
}

func Init(opt *InitOptions) {
	if opt == nil {
		opt = &InitOptions{}
	}

	if len(opt.ConfigFile) == 0 {
		Config = configs.NewConfig("data/config/conf.yaml")
	} else {
		Config = configs.NewConfig(opt.ConfigFile)
	}

	err := utils.CreateMutiDir(Config.DataPath)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", Config.DataPath)
	}
	err = utils.CreateMutiDir(Config.SavePath)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", Config.SavePath)
	}
	err = utils.CreateMutiDir(Config.CachePath)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", Config.CachePath)
	}

	if opt.Cache == nil {
		Cache = cache.NewBolt()
		Cache.Open(Config.Setting.CachePath)
	} else {
		Cache = opt.Cache
	}
}

package store

import (
	"GoBangumi/configs"
	"GoBangumi/internal/cache"
)

var (
	Cache  cache.Cache
	Config *configs.Config
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
		Config = configs.NewConfig("/Users/wetor/GoProjects/GoBangumi/data/config/conf.yaml")
	} else {
		Config = configs.NewConfig(opt.ConfigFile)
	}

	if opt.Cache == nil {
		Cache = cache.NewBolt()
		Cache.Open(Config.Setting.CachePath)
	} else {
		Cache = opt.Cache
	}
}

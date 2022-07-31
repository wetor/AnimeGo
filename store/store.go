package store

import (
	"GoBangumi/store/cache"
	"GoBangumi/store/config"
)

var (
	Cache  cache.Cache
	Config *config.Config
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
		Config = config.NewConfig("/Users/wetor/GoProjects/GoBangumi/data/config/conf.yaml")
	} else {
		Config = config.NewConfig(opt.ConfigFile)
	}

	if opt.Cache == nil {
		Cache = cache.NewBolt()
		Cache.Open(Config.Setting.CachePath)
	} else {
		Cache = opt.Cache
	}
}

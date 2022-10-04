package store

import (
	"AnimeGo/configs"
	"AnimeGo/internal/cache"
	"AnimeGo/internal/utils"
	"path"
	"sync"

	"go.uber.org/zap"
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
	InitDir()
	CreateDir()
	if opt.Cache == nil {
		Cache = cache.NewBolt()
		Cache.Open(Config.Setting.DbFile)
	} else {
		Cache = opt.Cache
	}
}

func InitDir() {
	Config.JavaScript = path.Join(Config.DataPath, Config.JavaScript)
	Config.DbFile = path.Join(Config.DataPath, Config.DbFile)
	CacheDir = path.Join(Config.DataPath, CacheDir)
}

func CreateDir() {
	err := utils.CreateMutiDir(Config.DataPath)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", Config.DataPath)
	}
	err = utils.CreateMutiDir(Config.SavePath)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", Config.SavePath)
	}
	err = utils.CreateMutiDir(CacheDir)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", CacheDir)
	}
	dbDir := path.Join(Config.DataPath, path.Dir(Config.DbFile))
	err = utils.CreateMutiDir(dbDir)
	if err != nil {
		zap.S().Fatalf("创建文件夹失败，%s", dbDir)
	}
}

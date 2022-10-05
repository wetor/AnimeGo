package store

import (
	"AnimeGo/configs"
	"AnimeGo/internal/cache"
	"AnimeGo/internal/process"
	"sync"
)

var (
	Cache   cache.Cache
	Config  *configs.Config
	WG      sync.WaitGroup
	Process process.Process
)

type InitOptions struct {
	Config *configs.Config
	Cache  cache.Cache
}

// Init
//  @Description: 初始化store和dir
//  @param opt *InitOptions
//
func Init(opt *InitOptions) {

	Config = opt.Config
	Cache = opt.Cache

}

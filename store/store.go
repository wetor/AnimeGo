package store

import (
	"GoBangumi/modules/cache"
)

const (
	InitStart = iota
	InitLoadConfig
	InitLoadCache
	InitConnectClient

	InitFinish
)

var (
	Cache     cache.Cache
	InitState int = InitStart
)

func SetCache(c cache.Cache) {
	Cache = c
}

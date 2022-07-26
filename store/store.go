package store

import (
	"GoBangumi/modules/cache"
)

var (
	Cache cache.Cache
)

func SetCache(c cache.Cache) {
	Cache = c
}

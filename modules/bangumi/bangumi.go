package bangumi

import (
	"GoBangumi/models"
	"GoBangumi/modules/cache"
)

type Bangumi interface {
	Parse(opt *models.BangumiParseOptions) *models.Bangumi
}

var Cache cache.Cache

func Init(c cache.Cache) {
	Cache = c
}

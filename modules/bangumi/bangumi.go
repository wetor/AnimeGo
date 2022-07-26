package bangumi

import (
	"GoBangumi/models"
	"GoBangumi/modules/cache"
)

var Cache cache.Cache

func Init(c cache.Cache) {
	Cache = c
}

type Bangumi interface {
	Parse(opt *models.BangumiParseOptions) *models.Bangumi
}

// Package anisource
// @Description: 番剧源包，用来收集番剧信息
package anisource

import (
	"AnimeGo/internal/cache"
	"AnimeGo/internal/models"
	"AnimeGo/pkg/anisource"
	"AnimeGo/pkg/anisource/bangumi"
	"AnimeGo/pkg/anisource/mikan"
	"AnimeGo/pkg/anisource/themoviedb"
)

type AniSource interface {
	Parse(opt *models.AnimeParseOptions) *models.AnimeEntity
}

// 单例模式
var (
	mikanInstance      *mikan.Mikan
	bangumiInstance    *bangumi.Bangumi
	themoviedbInstance *themoviedb.Themoviedb
)

// Init
//  @Description: 初始化anisource，需要在程序启动时调用
//  @param cache cache.Cache
//  @param proxy string
//
func Init(cache cache.Cache, proxy string) {
	mikanInstance = nil
	bangumiInstance = nil
	themoviedbInstance = nil
	anisource.Init(anisource.Options{
		Cache: cache,
		Proxy: proxy,
	})
}

func Mikan() *mikan.Mikan {
	if mikanInstance == nil {
		mikanInstance = &mikan.Mikan{}
	}
	return mikanInstance
}

func Bangumi() *bangumi.Bangumi {
	if bangumiInstance == nil {
		bangumiInstance = &bangumi.Bangumi{}
	}
	return bangumiInstance
}

func Themoviedb(key string) *themoviedb.Themoviedb {
	if themoviedbInstance == nil {
		themoviedbInstance = &themoviedb.Themoviedb{
			Key: key,
		}
	}
	return themoviedbInstance
}

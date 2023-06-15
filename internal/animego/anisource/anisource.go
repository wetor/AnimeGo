// Package anisource
// @Description: 番剧源包，用来收集番剧信息
package anisource

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/api"
)

// 单例模式
var (
	mikanInstance      api.AniDataParse
	bangumiInstance    api.AniDataGet
	themoviedbInstance api.AniDataSearchGet
)

type Options struct {
	*anidata.Options
}

// Init
//
//	@Description: 初始化anisource，需要在程序启动时调用
//	@param proxy string
func Init(opts *Options) {
	mikanInstance = nil
	bangumiInstance = nil
	themoviedbInstance = nil
	anidata.Init(opts.Options)
}

func Mikan() api.AniDataParse {
	if mikanInstance == nil {
		mikanInstance = &mikan.Mikan{}
	}
	return mikanInstance
}

func Bangumi() api.AniDataGet {
	if bangumiInstance == nil {
		bangumiInstance = &bangumi.Bangumi{}
	}
	return bangumiInstance
}

func Themoviedb(key string) api.AniDataSearchGet {
	if themoviedbInstance == nil {
		themoviedbInstance = &themoviedb.Themoviedb{
			Key: key,
		}
	}
	return themoviedbInstance
}

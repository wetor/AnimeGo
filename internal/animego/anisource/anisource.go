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
	MikanInstance      api.AniDataParse
	BangumiInstance    api.AniDataGet
	ThemoviedbInstance api.AniDataSearchGet
)

type Options struct {
	*anidata.Options
}

// Init
//
//	@Description: 初始化anisource，需要在程序启动时调用
//	@param proxy string
func Init(opts *Options) {
	MikanInstance = nil
	BangumiInstance = nil
	ThemoviedbInstance = nil
	anidata.Init(opts.Options)
}

func Mikan() api.AniDataParse {
	if MikanInstance == nil {
		MikanInstance = &mikan.Mikan{}
	}
	return MikanInstance
}

func Bangumi() api.AniDataGet {
	if BangumiInstance == nil {
		BangumiInstance = &bangumi.Bangumi{}
	}
	return BangumiInstance
}

func Themoviedb(key string) api.AniDataSearchGet {
	if ThemoviedbInstance == nil {
		ThemoviedbInstance = &themoviedb.Themoviedb{
			Key: key,
		}
	}
	return ThemoviedbInstance
}

package api

import "github.com/wetor/AnimeGo/internal/models"

type MikanInfo interface {
	CacheParseMikanInfo(url string) (mikanInfo any, err error)
}

type AniSource interface {
	Parse(opt *models.AnimeParseOptions) (*models.AnimeEntity, error)
}

type AniDataParse interface {
	Parse(options any) (any, error)
	ParseCache(options any) (any, error)
}

type AniDataSearch interface {
	Search(name string, filters any) (int, error)
	SearchCache(name string, filters any) (int, error)
}

type AniDataGet interface {
	Get(id int, filters any) (any, error)
	GetCache(id int, filters any) (any, error)
}

type AniDataParseSearch interface {
	AniDataParse
	AniDataSearch
}

type AniDataParseGet interface {
	AniDataParse
	AniDataGet
}

type AniDataSearchGet interface {
	AniDataSearch
	AniDataGet
}

type AniData interface {
	IName
}

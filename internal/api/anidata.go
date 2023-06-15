package api

type AniDataParse interface {
	Parse(options any) any
	ParseCache(options any) any
}

type AniDataSearch interface {
	Search(name string) int
	SearchCache(name string) int
}

type AniDataGet interface {
	Get(id int, filters any) any
	GetCache(id int, filters any) any
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
	AniDataParse
	AniDataSearch
	AniDataGet
}

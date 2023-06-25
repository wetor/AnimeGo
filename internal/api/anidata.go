package api

type AniDataParse interface {
	Parse(options any) (any, error)
	ParseCache(options any) (any, error)
}

type AniDataSearch interface {
	Search(name string) (int, error)
	SearchCache(name string) (int, error)
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

package lib

import "github.com/wetor/AnimeGo/internal/api"

var isInit = false
var (
	Feed api.Feed
)

type Options struct {
	Feed api.Feed
}

func Init(opts *Options) {
	if !isInit {
		InitLog()
		InitAnimeGo()
		isInit = true
		Feed = opts.Feed
	}
}

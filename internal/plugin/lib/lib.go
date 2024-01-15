package lib

import "github.com/wetor/AnimeGo/internal/api"

var isInit = false
var (
	Feed  api.Feed
	Mikan api.MikanInfo
)

type Options struct {
	Feed  api.Feed
	Mikan api.MikanInfo
}

func Init(opts *Options) {
	if !isInit {
		InitLog()
		InitAnimeGo()
		isInit = true
	}
	Feed = opts.Feed
	Mikan = opts.Mikan
}

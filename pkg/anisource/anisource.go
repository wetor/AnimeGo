package anisource

import mem "github.com/wetor/AnimeGo/pkg/memorizer"

var (
	Cache mem.Memorizer
)

type Options struct {
	Cache mem.Memorizer
}

func Init(opts *Options) {
	Cache = opts.Cache
}

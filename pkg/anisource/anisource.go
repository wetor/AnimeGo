package anisource

import mem "AnimeGo/pkg/memorizer"

var (
	Cache mem.Memorizer
)

type Options struct {
	Cache mem.Memorizer
}

func Init(opts *Options) {
	Cache = opts.Cache
}

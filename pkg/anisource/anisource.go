package anisource

import mem "github.com/wetor/AnimeGo/pkg/memorizer"

var (
	Cache     mem.Memorizer
	CacheTime map[string]int64
)

type Options struct {
	Cache     mem.Memorizer
	CacheTime map[string]int64
}

func Init(opts *Options) {
	Cache = opts.Cache
	CacheTime = opts.CacheTime
}

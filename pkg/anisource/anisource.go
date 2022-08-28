package anisource

import mem "AnimeGo/pkg/memorizer"

var (
	Proxy   string
	Timeout int
	Retry   int
	Cache   mem.Memorizer
)

type Options struct {
	Proxy   string
	Timeout int
	Retry   int
	Cache   mem.Memorizer
}

func Init(opts *Options) {
	Proxy = opts.Proxy
	Timeout = opts.Timeout
	Retry = opts.Retry
	Cache = opts.Cache
}

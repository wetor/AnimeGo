package anisource

import mem "GoBangumi/pkg/memorizer"

var (
	Proxy string
	Cache mem.Memorizer = nil
)

type Options struct {
	Proxy string
	Cache mem.Memorizer
}

func Init(opts Options) {
	Proxy = opts.Proxy
	Cache = opts.Cache
}

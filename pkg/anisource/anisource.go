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

func (o *Options) Default() {
	if o.Timeout == 0 {
		o.Timeout = 3
	}
	if o.Retry == 0 {
		o.Retry = 1
	}
}

func Init(opts *Options) {
	opts.Default()
	Proxy = opts.Proxy
	Timeout = opts.Timeout
	Retry = opts.Retry
	Cache = opts.Cache
}

package downloader

import "sync"

var (
	RefreshSecond int
	Category      string
	WG            *sync.WaitGroup
	ReInitWG      sync.WaitGroup
)

type Options struct {
	RefreshSecond int
	Category      string
	WG            *sync.WaitGroup
}

func Init(opts *Options) {
	RefreshSecond = opts.RefreshSecond
	Category = opts.Category
	WG = opts.WG
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	RefreshSecond = opts.RefreshSecond
	Category = opts.Category
}

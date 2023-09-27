package renamer

import "sync"

var (
	WG            *sync.WaitGroup
	RefreshSecond int
	ReInitWG      sync.WaitGroup
)

type Options struct {
	WG            *sync.WaitGroup
	RefreshSecond int
}

func Init(opts *Options) {
	WG = opts.WG
	RefreshSecond = opts.RefreshSecond
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
}

package renamer

import "sync"

var (
	WG                *sync.WaitGroup
	UpdateDelaySecond int
	ReInitWG          sync.WaitGroup
)

type Options struct {
	WG                *sync.WaitGroup
	UpdateDelaySecond int
}

func Init(opts *Options) {
	WG = opts.WG
	UpdateDelaySecond = opts.UpdateDelaySecond
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
}

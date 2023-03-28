package filter

import "sync"

var (
	DelaySecond int
	ReInitWG    sync.WaitGroup
)

type Options struct {
	DelaySecond int
}

func Init(opts *Options) {
	DelaySecond = opts.DelaySecond
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	DelaySecond = opts.DelaySecond
}

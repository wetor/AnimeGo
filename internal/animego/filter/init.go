package filter

import "sync"

var (
	MultiGoroutineMax     int
	MultiGoroutineEnabled bool
	DelaySecond           int
	ReInitWG              sync.WaitGroup
)

type Options struct {
	MultiGoroutineMax     int
	MultiGoroutineEnabled bool
	DelaySecond           int
}

func Init(opts *Options) {
	MultiGoroutineEnabled = opts.MultiGoroutineEnabled
	MultiGoroutineMax = opts.MultiGoroutineMax
	DelaySecond = opts.DelaySecond

}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	MultiGoroutineEnabled = opts.MultiGoroutineEnabled
	MultiGoroutineMax = opts.MultiGoroutineMax
	DelaySecond = opts.DelaySecond
}

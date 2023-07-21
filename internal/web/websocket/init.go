package websocket

import "sync"

var (
	WG       *sync.WaitGroup
	Notify   chan []byte
	ReInitWG sync.WaitGroup
)

type Options struct {
	WG     *sync.WaitGroup
	Notify chan []byte
}

func Init(opts *Options) {
	WG = opts.WG
	Notify = opts.Notify
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	Notify = opts.Notify
}

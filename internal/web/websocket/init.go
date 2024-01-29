package websocket

import "sync"

var (
	WG     *sync.WaitGroup
	Notify chan []byte
)

type Options struct {
	WG     *sync.WaitGroup
	Notify chan []byte
}

func Init(opts *Options) {
	WG = opts.WG
	Notify = opts.Notify
}

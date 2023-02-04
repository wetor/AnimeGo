package logger

import (
	"context"
	"sync"

	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Options struct {
	File    string
	Debug   bool
	Context context.Context
	WG      *sync.WaitGroup
}

func Init(opts *Options) {
	log.Init(&log.Options{
		File:  opts.File,
		Debug: opts.Debug,
	})
	opts.WG.Add(1)
	go func() {
		defer opts.WG.Done()
		for {
			select {
			case <-opts.Context.Done():
				Flush()
				log.Debugf("正常退出 logger")
				return
			default:
				Flush()
				utils.Sleep(30, opts.Context)
			}
		}
	}()
}

func Flush() {
	_ = log.Sync()
}

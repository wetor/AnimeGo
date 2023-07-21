package logger

import (
	"context"
	"io"
	"sync"

	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const logSyncSecond = 30

type Options struct {
	File    string
	Debug   bool
	Context context.Context
	WG      *sync.WaitGroup
	Out     io.Writer
}

func Init(opts *Options) {
	log.Init(&log.Options{
		File:  opts.File,
		Debug: opts.Debug,
		Out:   opts.Out,
	})
	opts.WG.Add(1)
	go func() {
		defer opts.WG.Done()
		for {
			select {
			case <-opts.Context.Done():
				_ = log.Sync()
				_ = log.Close()
				log.Debugf("正常退出 logger")
				return
			default:
				_ = log.Sync()
				utils.Sleep(logSyncSecond, opts.Context)
			}
		}
	}()
}

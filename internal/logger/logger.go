package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/utils"
)

type Options struct {
	File    string
	Debug   bool
	Context context.Context
	WG      *sync.WaitGroup
}

func Init(opts *Options) {
	opts.WG.Add(1)

	GetLogger(opts)
	go func() {
		defer opts.WG.Done()
		for {
			select {
			case <-opts.Context.Done():
				Flush()
				zap.S().Debug("正常退出 logger")
				return
			default:
				Flush()
				utils.Sleep(30, opts.Context)
			}
		}
	}()
}

func Flush() {
	_ = zap.S().Sync()
}

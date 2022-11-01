package logger

import (
	"context"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"go.uber.org/zap"
)

type InitOptions struct {
	File    string
	Debug   bool
	Context context.Context
}

func Init(opt *InitOptions) {
	store.WG.Add(1)

	GetLogger(opt)
	go func() {
		defer store.WG.Done()
		for {
			select {
			case <-opt.Context.Done():
				Flush()
				zap.S().Debug("正常退出 logger")
				return
			default:
				Flush()
				utils.Sleep(30, opt.Context)
			}
		}
	}()
}

func Flush() {
	_ = zap.S().Sync()
}

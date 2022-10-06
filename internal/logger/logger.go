package logger

import (
	"AnimeGo/internal/store"
	"AnimeGo/internal/utils"
	"context"
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
				zap.S().Info("正常退出")
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

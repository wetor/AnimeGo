package logger

import (
	"go.uber.org/zap"
	"time"
)

func Init() {
	GetLogger()
	go func() {
		Flush()
		time.Sleep(30 * time.Second)
	}()
}

func Flush() {
	_ = zap.S().Sync()
}

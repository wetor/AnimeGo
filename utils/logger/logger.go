package logger

import "go.uber.org/zap"

func Init() {
	GetLogger()
}

func Flush() {
	_ = zap.S().Sync()
}

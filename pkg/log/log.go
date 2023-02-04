package log

import (
	"os"
	"path"
	"sync"

	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger
	loggerInit sync.Once
	file       = "log/log.log"
	debug      = false
)

type Options struct {
	File  string
	Debug bool
}

func Init(opts *Options) {
	file = opts.File
	debug = opts.Debug
	dir := path.Dir(file)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			_ = os.MkdirAll(dir, os.ModePerm)
		}
	}
	logger = NewLogger(file, debug).Sugar()
}

func GetLogger() *zap.SugaredLogger {
	loggerInit.Do(func() {
		if logger == nil {
			// logger is not initialized, for example, running `go test`
			logger = NewLogger(file, debug).Sugar()
		}
	})
	return logger
}

func Sync() error {
	return GetLogger().Sync()
}

func Infof(template string, args ...interface{}) {
	GetLogger().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	GetLogger().Warnf(template, args...)
}

func Debugf(template string, args ...interface{}) {
	GetLogger().Debugf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	GetLogger().Errorf(template, args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}

func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

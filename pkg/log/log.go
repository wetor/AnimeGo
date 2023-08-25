package log

import (
	"io"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var (
	logger     *zap.SugaredLogger
	loggerInit sync.Once
	file       = "log/log.log"
	debug      = false
	out        io.Writer
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type Options struct {
	File  string
	Debug bool
	Out   io.Writer
}

func Init(opts *Options) {
	debug = opts.Debug
	out = opts.Out
	file = opts.File
	dir := path.Dir(file)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsExist(err) {
			_ = os.MkdirAll(dir, os.ModePerm)
		}
	}
	logger = NewLogger(file, debug, out).Sugar()
}

func ReInt(opts *Options) {
	if len(opts.File) != 0 {
		file = opts.File
	}
	if opts.Out != nil {
		out = opts.Out
	}
	debug = opts.Debug
	logger = NewLogger(file, debug, out).Sugar()
}

func GetLogger() *zap.SugaredLogger {
	loggerInit.Do(func() {
		if logger == nil {
			// logger is not initialized, for example, running `go test`
			logger = NewLogger(file, debug, out).Sugar()
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

func DebugErr(err error) {
	if e, ok := err.(StackTracer); ok {
		st := e.StackTrace()
		GetLogger().Debugf("%s%+v", err, st[0:2])
	} else {
		GetLogger().Debugf("%+v", err)
	}
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

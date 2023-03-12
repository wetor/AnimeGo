package lib

import (
	"github.com/go-python/gpython/py"
	log2 "github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/plugin"
)

func InitLog() {
	methods := []*py.Method{
		py.MustNewMethod("debug", log("debug"), 0, "debug(args...)"),
		py.MustNewMethod("debugf", log("debugf"), 0, "debugf(template, args...)"),
		py.MustNewMethod("info", log("info"), 0, "info(args...)"),
		py.MustNewMethod("infof", log("infof"), 0, "infof(template, args...)"),
		py.MustNewMethod("error", log("error"), 0, "error(args...)"),
		py.MustNewMethod("errorf", log("errorf"), 0, "errorf(template, args...)"),
		py.MustNewMethod("warn", log("warn"), 0, "warn(args...)"),
		py.MustNewMethod("warnf", log("warnf"), 0, "warnf(template, args...)"),
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "log",
			Doc:  "Log Module",
		},
		Methods: methods,
	})
}

func log(name string) func(self py.Object, args py.Tuple) (py.Object, error) {
	var logFunc func(args ...interface{})
	var logFuncf func(template string, args ...interface{})
	switch name {
	case "debug":
		logFunc = log2.Debug
	case "info":
		logFunc = log2.Info
	case "warn":
		logFunc = log2.Warn
	case "error":
		logFunc = log2.Error
	case "debugf":
		logFuncf = log2.Debugf
	case "infof":
		logFuncf = log2.Infof
	case "warnf":
		logFuncf = log2.Warnf
	case "errorf":
		logFuncf = log2.Errorf
	}
	if name[len(name)-1] == 'f' {
		return func(self py.Object, args py.Tuple) (py.Object, error) {
			if len(args) < 1 {
				return py.AttributeError, nil
			}
			f := plugin.PyObject2Value(args[0]).(string)
			list := make([]any, len(args)-1)
			for i := 0; i < len(args)-1; i++ {
				list[i] = plugin.PyObject2Value(args[i+1])
			}
			logFuncf(f, list...)

			return py.None, nil
		}
	} else {
		return func(self py.Object, args py.Tuple) (py.Object, error) {
			list := make([]any, len(args))
			for i := 0; i < len(args); i++ {
				list[i] = plugin.PyObject2Value(args[i])
			}
			logFunc(list...)

			return py.None, nil
		}
	}
}

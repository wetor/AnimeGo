package lib

import (
	"github.com/go-python/gpython/py"
	pyutils "github.com/wetor/AnimeGo/internal/plugin/python/utils"
	"go.uber.org/zap"
)

func InitLog() {
	methods := []*py.Method{
		py.MustNewMethod("debug", log("debug"), 0, "debug(args...)"),
		py.MustNewMethod("debugf", log("debugf"), 0, "debugf(args...)"),
		py.MustNewMethod("info", log("info"), 0, "info(args...)"),
		py.MustNewMethod("infof", log("infof"), 0, "infof(args...)"),
		py.MustNewMethod("error", log("error"), 0, "error(args...)"),
		py.MustNewMethod("errorf", log("errorf"), 0, "errorf(args...)"),
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
		logFunc = zap.S().Debug
	case "info":
		logFunc = zap.S().Info
	case "error":
		logFunc = zap.S().Error
	case "debugf":
		logFuncf = zap.S().Debugf
	case "infof":
		logFuncf = zap.S().Infof
	case "errorf":
		logFuncf = zap.S().Errorf
	}
	if name[len(name)-1] == 'f' {
		return func(self py.Object, args py.Tuple) (py.Object, error) {
			if len(args) < 1 {
				return py.AttributeError, nil
			}
			f := pyutils.PyObject2Value(args[0]).(string)
			list := make([]any, len(args)-1)
			for i := 0; i < len(args)-1; i++ {
				list[i] = pyutils.PyObject2Value(args[i+1])
			}
			logFuncf(f, list...)

			return py.None, nil
		}
	} else {
		return func(self py.Object, args py.Tuple) (py.Object, error) {
			list := make([]any, len(args))
			for i := 0; i < len(args); i++ {
				list[i] = pyutils.PyObject2Value(args[i])
			}
			logFunc(list...)

			return py.None, nil
		}
	}
}
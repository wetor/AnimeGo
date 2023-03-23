package lib

import (
	"github.com/go-python/gpython/py"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
)

func InitLog() {
	methods := []*py.Method{
		py.MustNewMethod("debug", logBase("debug"), 0, "debug(args...)"),
		py.MustNewMethod("debugf", logBase("debugf"), 0, "debugf(template, args...)"),
		py.MustNewMethod("info", logBase("info"), 0, "info(args...)"),
		py.MustNewMethod("infof", logBase("infof"), 0, "infof(template, args...)"),
		py.MustNewMethod("error", logBase("error"), 0, "error(args...)"),
		py.MustNewMethod("errorf", logBase("errorf"), 0, "errorf(template, args...)"),
		py.MustNewMethod("warn", logBase("warn"), 0, "warn(args...)"),
		py.MustNewMethod("warnf", logBase("warnf"), 0, "warnf(template, args...)"),
	}

	py.RegisterModule(&py.ModuleImpl{
		Info: py.ModuleInfo{
			Name: "log",
			Doc:  "Log Module",
		},
		Methods: methods,
	})
}

func debugMode(self py.Object) bool {
	m, _ := self.(*py.Module).Context.GetModule("__main__")
	if m != nil {
		if d, ok := m.Globals["__debug__"]; ok {
			if debug, ok := d.(py.Bool); ok {
				if debug {
					return true
				}
			}
		}
	}
	// 未找到__main__模块，未设置__debug__，已设置但不是bool类型，已设置但结果为false
	return false
}

func logBase(name string) func(self py.Object, args py.Tuple) (py.Object, error) {
	var logFunc func(args ...interface{})
	var logFuncf func(template string, args ...interface{})
	switch name {
	case "debug":
		logFunc = log.Debug
	case "info":
		logFunc = log.Info
	case "warn":
		logFunc = log.Warn
	case "error":
		logFunc = log.Error
	case "debugf":
		logFuncf = log.Debugf
	case "infof":
		logFuncf = log.Infof
	case "warnf":
		logFuncf = log.Warnf
	case "errorf":
		logFuncf = log.Errorf
	}
	if name[len(name)-1] == 'f' {
		return func(self py.Object, args py.Tuple) (py.Object, error) {
			if len(args) < 1 {
				return py.AttributeError, nil
			}
			if name == "debugf" {
				if !debugMode(self) {
					return py.None, nil
				}
			}
			f := python.ToValue(args[0]).(string)
			list := make([]any, len(args)-1)
			for i := 0; i < len(args)-1; i++ {
				list[i] = python.ToValue(args[i+1])
			}
			logFuncf(f, list...)

			return py.None, nil
		}
	} else {
		return func(self py.Object, args py.Tuple) (py.Object, error) {
			if name == "debug" {
				if !debugMode(self) {
					return py.None, nil
				}
			}
			list := make([]any, len(args))
			for i := 0; i < len(args); i++ {
				list[i] = python.ToValue(args[i])
			}
			logFunc(list...)

			return py.None, nil
		}
	}
}

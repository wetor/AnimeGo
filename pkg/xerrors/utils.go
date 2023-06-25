package xerrors

import (
	"runtime"
)

func GetCaller(skip int) (string, string, int) {
	pc, file, line, _ := runtime.Caller(skip + 1)
	pcName := runtime.FuncForPC(pc).Name()
	return file, pcName, line
}

func HandleError(fn func(error)) {
	if err := recover(); err != nil {
		if e, ok := err.(error); ok {
			fn(e)
		} else {
			panic(err)
		}
	}
}

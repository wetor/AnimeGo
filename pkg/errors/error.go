package errors

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
)

type AniError struct {
	Msg  string
	File string //file
	Func string //func
	Line int
}

func NewAniError(msg string) *AniError {
	return NewAniErrorSkipf(2, msg)
}

func NewAniErrorSkip(skip int, msg string) *AniError {
	return NewAniErrorSkipf(skip+1, msg)
}

func NewAniErrorf(format string, a ...interface{}) *AniError {
	return NewAniErrorSkipf(2, format, a...)
}

func NewAniErrorSkipf(skip int, format string, a ...interface{}) *AniError {
	pc, file, line, _ := runtime.Caller(skip)
	pcName := runtime.FuncForPC(pc).Name()
	return &AniError{
		Msg:  fmt.Sprintf(format, a...),
		File: file,
		Func: pcName,
		Line: line,
	}
}

func (e *AniError) Error() string {
	str := bytes.NewBuffer(nil)
	str.WriteString("[Msg]: ")
	str.WriteString(e.Msg)
	if len(e.Func) > 0 {
		str.WriteString(", [Func]: ")
		str.WriteString(e.Func)
	}
	if len(e.File) > 0 {
		_, file := path.Split(e.File)
		str.WriteString(fmt.Sprintf(", [File]: %s:%d", file, e.Line))
	}

	return str.String()
}

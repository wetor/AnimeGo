package errors

import (
	"bytes"
	"fmt"
	"path"
	"runtime"
)

type AniError struct {
	Data interface{}
	Msg  string
	File string //file
	Func string //func
	Line int
}

func NewAniError(msg string) *AniError {
	return NewAniErrorSkipf(2, msg, nil)
}

func NewAniErrorD(data interface{}) *AniError {
	return NewAniErrorSkipf(2, "", data)
}

func NewAniErrorSkip(skip int, msg string) *AniError {
	return NewAniErrorSkipf(skip+1, msg, nil)
}

func NewAniErrorf(format string, a ...interface{}) *AniError {
	return NewAniErrorSkipf(2, format, nil, a...)
}

func NewAniErrorSkipf(skip int, format string, data interface{}, a ...interface{}) *AniError {
	pc, file, line, _ := runtime.Caller(skip)
	pcName := runtime.FuncForPC(pc).Name()
	return &AniError{
		Msg:  fmt.Sprintf(format, a...),
		Data: data,
		File: file,
		Func: pcName,
		Line: line,
	}
}

func (e *AniError) SetMsg(msg string) *AniError {
	e.Msg = msg
	return e
}
func (e *AniError) SetData(data interface{}) *AniError {
	e.Data = data
	return e
}

func (e *AniError) Error() string {
	str := bytes.NewBuffer(nil)
	str.WriteString(e.Msg)

	_, file := path.Split(e.File)
	str.WriteString(fmt.Sprintf(" [(%s) %s:%d]", e.Func, file, e.Line))

	return str.String()
}

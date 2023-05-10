package logger

import (
	"bytes"
	"io"
	"os"
)

var enableNotify = false

type LogNotify struct {
	notify chan []byte
}

func (w *LogNotify) Write(p []byte) (n int, err error) {
	n, err = os.Stdout.Write(p)
	if enableNotify && err == nil {
		b := bytes.Clone(p)
		go func(data []byte) {
			w.notify <- data
		}(b)
	}
	return
}

func NewLogNotify() (io.Writer, chan []byte) {
	notify := make(chan []byte)
	return &LogNotify{
		notify: notify,
	}, notify
}

func SetLogNotify(enable bool) {
	enableNotify = enable
}

func GetLogNotify() bool {
	return enableNotify
}

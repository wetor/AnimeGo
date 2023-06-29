package logger

import (
	"bytes"
	"io"
	"os"
)

const (
	NotifyDisabled = iota
	NotifyEnabled
	NotifyPaused
)

var notifyStatus = NotifyDisabled

type LogNotify struct {
	notify chan []byte
}

func (w *LogNotify) Write(p []byte) (n int, err error) {
	n, err = os.Stdout.Write(p)
	if notifyStatus >= NotifyEnabled && err == nil {
		b := bytes.Clone(p)
		w.notify <- b
	}
	return
}

func NewLogNotify() (io.Writer, chan []byte) {
	notify := make(chan []byte, 30)
	return &LogNotify{
		notify: notify,
	}, notify
}

func DisableLogNotify() {
	notifyStatus = NotifyDisabled
}

func EnableLogNotify() {
	notifyStatus = NotifyEnabled
}

func PauseLogNotify() {
	notifyStatus = NotifyPaused
}

func GetLogNotify() int {
	return notifyStatus
}

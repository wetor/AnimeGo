package logger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/pkg/log"
)

func TestNewLogNotify(t *testing.T) {
	out, notify := logger.NewLogNotify()
	logger.SetLogNotify(true)
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
		Out:   out,
	})

	go func() {
		for {
			select {
			case data := <-notify:
				fmt.Printf("Recive: %s", data)
			}
		}
	}()

	time.Sleep(1 * time.Second)
	log.Debugf("test %d, 测试文本", 10086)
	//time.Sleep(100 * time.Microsecond)
	log.Warnf("test, warn")
	//time.Sleep(100 * time.Microsecond)
	log.Infof("test, info1")
	time.Sleep(1 * time.Second)
	log.Infof("test, info2")
	time.Sleep(1 * time.Second)
}

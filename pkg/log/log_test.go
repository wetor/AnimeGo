package log_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/pkg/log"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	m.Run()
	fmt.Println("end")
}

func TestDebug(t *testing.T) {
	log.Debugf("%s %s", "测试文本", "111")
	log.Infof("222")
}

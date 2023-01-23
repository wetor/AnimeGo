package request

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	m.Run()
	fmt.Println("end")
}

func TestGet(t *testing.T) {
	Init(&Options{})
	str, err := GetString("http://pv.sohu.com/cityjson?ie=utf-8")
	fmt.Println(str, err)
}

func TestGetRetry(t *testing.T) {
	Init(&Options{
		Retry:     1,
		RetryWait: 0,
		Timeout:   0,
		Debug:     true,
	})
	str, err := GetString("https://www.baidu.com/aaa/test")
	fmt.Println(str, err)
}

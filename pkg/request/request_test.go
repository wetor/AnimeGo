package request

import (
	"fmt"
	"go.uber.org/zap"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	m.Run()
	fmt.Println("end")
}

func TestGet(t *testing.T) {
	Init(&InitOptions{})

	str, err := GetString("http://pv.sohu.com/cityjson?ie=utf-8")
	fmt.Println(str, err.Error())
}

func TestGetRetry(t *testing.T) {
	Init(&InitOptions{
		Proxy:     "http://127.0.0.1:7890",
		Retry:     1,
		RetryWait: 0,
		Timeout:   1,
		Debug:     true,
	})
	str, err := GetString("https://www.baidu.com/aaa/test")
	fmt.Println(str, err.Error())
}

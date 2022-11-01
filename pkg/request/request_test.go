package request

import (
	"fmt"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestGet(t *testing.T) {
	Init(&InitOptions{})

	err, str := GetString("http://pv.sohu.com/cityjson?ie=utf-8")
	if err != nil {
		panic(err)
	}
	fmt.Println(str)
}

func TestGetRetry(t *testing.T) {
	test.TestInit()
	Init(&InitOptions{
		Proxy:     "http://127.0.0.1:7890",
		Retry:     3,
		RetryWait: 1,
		Timeout:   3,
		Debug:     true,
	})
	err, str := GetString("https://www.baidu.com/aaa/test")
	fmt.Println(str)
	if err != nil {
		panic(err)
	}
}

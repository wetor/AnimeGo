package request_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	m.Run()
	fmt.Println("end")
}

func TestGet(t *testing.T) {
	request.Init(&request.Options{})
	str, err := request.GetString("http://pv.sohu.com/cityjson?ie=utf-8")
	fmt.Println(str, err)
}

func TestGetRetry(t *testing.T) {
	request.Init(&request.Options{
		Retry:     1,
		RetryWait: 0,
		Timeout:   0,
		Debug:     true,
	})
	str, err := request.GetString("https://www.baidu.com/aaa/test")
	fmt.Println(str, err)
}

package request_test

import (
	"bytes"
	"fmt"
	"os"
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
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestGet(t *testing.T) {
	request.Init(&request.Options{})
	rss := bytes.NewBuffer(nil)
	err := request.GetWriter("https://mikanani.me/RSS/MyBangumi?token=ky5DTt%2fMyAjCH2oKEN81FQ%3d%3d", rss)
	fmt.Println(rss.String(), err)
}

func TestGet2(t *testing.T) {
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

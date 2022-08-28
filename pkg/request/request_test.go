package request

import (
	"GoBangumi/internal/store"
	"GoBangumi/third_party/goreq"
	"fmt"
	"io"
	"testing"
)

func Test_goreq(t *testing.T) {
	req, err := goreq.Request{
		Method: "GET",
		Proxy:  "http://127.0.0.1:7890",
		Uri:    "http://pv.sohu.com/cityjson?ie=utf-8",
	}.Do()
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	fmt.Println(req.StatusCode)
	body, _ := io.ReadAll(req.Body)
	fmt.Println(string(body))
}

func TestGet(t *testing.T) {
	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/GoBangumi/data/config/conf.yaml",
	})
	err := Get(&Param{
		Uri: "http://pv.sohu.com/cityjson?ie=utf-8",
	})
	if err != nil {
		panic(err)
	}
}

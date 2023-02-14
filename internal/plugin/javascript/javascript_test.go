package javascript_test

import (
	"fmt"
	"testing"

	"github.com/dop251/goja"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/pkg/log"
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

func TestJs(t *testing.T) {
	vm := goja.New()
	_, err := vm.RunString(`
function sum(a, b) {
    return a+b;
}
`)
	if err != nil {
		panic(err)
	}
	sum, ok := goja.AssertFunction(vm.Get("sum"))
	if !ok {
		panic("Not a function")
	}

	res, err := sum(goja.Undefined(), vm.ToValue(40), vm.ToValue(2))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func TestJavaScript_Execute(t *testing.T) {
	js := &javascript.JavaScript{}
	js.SetSchema([]string{"feedItems,optional"}, []string{"index,optional", "error,optional"})
	execute := js.Execute(&models.PluginExecuteOptions{
		File: "testdata/test.js",
	}, models.Object{
		"feedItems": []*models.FeedItem{
			{
				Url:      "localhost",
				Name:     "【喵萌奶茶屋】★04月新番★[相合之物/Deaimon][09][720p][简体][招募翻译校对]",
				Download: "asdasdasd",
			},
			{
				Url:      "localhost:99",
				Name:     "[悠哈璃羽字幕社] [国王排名_Ousama Ranking] [22] [x264 1080p] [CHT]",
				Download: "asdasasdaddasd",
			},
		},
	})
	fmt.Println(execute)
}

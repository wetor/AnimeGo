package javascript

import (
	"AnimeGo/internal/models"
	"AnimeGo/test"
	"fmt"
	"github.com/dop251/goja"
	"testing"
)

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
	test.TestInit()

	js := &JavaScript{}
	js.SetSchema([]string{"feedItems"}, []string{"index", "error"})
	execute, err := js.Execute("/Users/wetor/GoProjects/AnimeGo/data/plugin/filter/default.js",
		Object{
			"feedItems": []*models.FeedItem{
				{
					Url:     "localhost",
					Name:    "【喵萌奶茶屋】★04月新番★[相合之物/Deaimon][09][720p][简体][招募翻译校对]",
					Torrent: "asdasdasd",
				},
				{
					Url:     "localhost:99",
					Name:    "[悠哈璃羽字幕社] [国王排名_Ousama Ranking] [22] [x264 1080p] [CHT]",
					Torrent: "asdasasdaddasd",
				},
			},
		})
	if err != nil {
		panic(err)
	}
	fmt.Println(execute)
}

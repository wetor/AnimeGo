package python

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/internal/plugin/python/lib"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	gpython.Init()
	m.Run()
	fmt.Println("end")
}

func TestPython_Execute(t *testing.T) {
	p := &Python{}
	p.SetSchema([]string{"title"}, []string{})
	result := p.Execute("data/raw_parser.py", models.Object{
		"title": "【悠哈璃羽字幕社】 [明日同学的水手服_Akebi-chan no Sailor-fuku] [01-12] [x264 1080p][CHT]",
	})
	fmt.Println(result)
}

func TestParser(t *testing.T) {
	fr, err := os.Open("./data/test_data.txt")
	if err != nil {
		panic(err)
	}
	defer fr.Close()
	sc := bufio.NewScanner(fr)

	var eps []models.Object
	p := &Python{}
	p.SetSchema([]string{"title"}, []string{})
	for sc.Scan() {
		title := sc.Text()
		fmt.Println(title)
		result := p.Execute("/Users/wetor/GoProjects/AnimeGo/data/plugin/anisource/Auto_Bangumi/raw_parser.py", models.Object{
			"title": title,
		})
		fmt.Println(result)
		eps = append(eps, result.(models.Object))
	}

	fw, err := os.Create("./data/test_out.txt")
	if err != nil {
		panic(err)
	}
	defer fw.Close()
	for _, ep := range eps {
		fw.WriteString(fmt.Sprintf("%v\n", ep))
	}
}

func TestLib_log(t *testing.T) {
	lib.InitLog()
	p := &Python{}
	p.SetSchema([]string{}, []string{})
	result := p.Execute("data/test_log.py", models.Object{
		"title": "【悠哈璃羽字幕社】 [明日同学的水手服_Akebi-chan no Sailor-fuku] [01-12] [x264 1080p][CHT]",
	})
	fmt.Println(result)
}

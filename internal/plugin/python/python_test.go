package python_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	gpython.Init()
	m.Run()
	fmt.Println("end")
}

func TestPython_Execute(t *testing.T) {
	p := &python.Python{}
	p.SetSchema([]string{"optional:title"}, []string{})
	result := p.Execute(&models.PluginExecuteOptions{
		File: "data/raw_parser.py",
	}, models.Object{
		"title": "[OPFans枫雪动漫][ONE PIECE 海贼王][第1048话][周日版][1080p][MP4][简体]",
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
	p := &python.Python{}
	p.SetSchema([]string{"optional:title"}, []string{})
	for sc.Scan() {
		title := sc.Text()
		fmt.Println(title)
		result := p.Execute(&models.PluginExecuteOptions{
			File: "/Users/wetor/GoProjects/AnimeGo/data/plugin/lib/Auto_Bangumi/raw_parser.py",
		}, models.Object{
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
	p := &python.Python{}
	p.SetSchema([]string{"optional:title"}, []string{"optional:result"})
	result := p.Execute(&models.PluginExecuteOptions{
		File: "data/test_log.py",
	}, models.Object{
		"title": "【悠哈璃羽字幕社】 [明日同学的水手服_Akebi-chan no Sailor-fuku] [01-12] [x264 1080p][CHT]",
	})
	fmt.Println(result)
}

func TestPython(t *testing.T) {
	lib.InitLog()
	p := &python.Python{}
	p.SetSchema([]string{"optional:title"}, []string{"optional:result"})
	result := p.Execute(&models.PluginExecuteOptions{
		File: "data/test.py",
	}, models.Object{})
	fmt.Println(result)
}

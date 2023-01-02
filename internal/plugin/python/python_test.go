package python

import (
	"bufio"
	"fmt"
	"github.com/go-python/gpython/py"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/third_party/gpython"
	"os"
	"strconv"
	"testing"
)

func TestRe1(t *testing.T) {
	i, _ := strconv.ParseInt("09", 10, 64)
	fmt.Println(i)
	gpython.Init()
	pyFile := "./data/test.py"
	// See type Context interface and related docs
	ctx := py.NewContext(py.DefaultContextOpts())

	// This drives modules being able to perform cleanup and release resources
	defer ctx.Close()

	_, err := py.RunFile(ctx, pyFile, py.CompileOpts{}, nil)

	if err != nil {
		py.TracebackDump(err)
		panic(err)
	}

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
	gpython.Init()
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

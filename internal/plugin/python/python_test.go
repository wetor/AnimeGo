package python_test

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/pkg/cache"

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
	p.Load(&models.PluginLoadOptions{
		File: "data/raw_parser.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:         "main",
				ParamsSchema: []string{"title"},
				ResultSchema: []string{"ep"},
			},
		},
	})
	result := p.Run("main", models.Object{
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

	p.Load(&models.PluginLoadOptions{
		File: "/Users/wetor/GoProjects/AnimeGo/data/plugin/lib/Auto_Bangumi/raw_parser.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})

	for sc.Scan() {
		title := sc.Text()
		fmt.Println(title)
		result := p.Run("main", models.Object{
			"title": title,
		})
		fmt.Println(result)
		eps = append(eps, result)
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
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "data/test_log.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:         "main",
				ParamsSchema: []string{"title"},
				ResultSchema: []string{"result"},
			},
		},
	})
	result := p.Run("main", models.Object{
		"title": "【悠哈璃羽字幕社】 [明日同学的水手服_Akebi-chan no Sailor-fuku] [01-12] [x264 1080p][CHT]",
	})
	fmt.Println(result)
}

func TestPython(t *testing.T) {
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "data/test.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
				ParamsSchema:    []string{"title"},
				ResultSchema:    []string{"result"},
			},
		},
	})
	result := p.Run("main", models.Object{})
	fmt.Println(result)
}

func TestPythonFunction(t *testing.T) {
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "testdata/test.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:         "main",
				ParamsSchema: []string{"params"},
				ResultSchema: []string{"result"},
			},
			{
				Name:            "test",
				SkipSchemaCheck: true,
			},
		},
	})
	result := p.Run("main", models.Object{
		"params": []int{1, 2, 3},
	})
	fmt.Println(result)

	p.Run("test", models.Object{
		"test": true,
	})
}

func TestPythonVariable(t *testing.T) {
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "testdata/var.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
		Variables: []*models.PluginVariableOptions{
			{
				Name: "Name",
			},
			{
				Name: "Cron",
			},
			{
				Name:     "Test",
				Nullable: true,
			},
		},
	})

	fmt.Println(p.Get("Name"))
	fmt.Println(p.Get("Cron"))
	p.Set("Name", "update_test")
	result := p.Run("main", models.Object{
		"params": []int{1, 2, 3},
	})
	fmt.Println(result)

}

func TestPythonJson(t *testing.T) {
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "testdata/json.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})
	result := p.Run("main", models.Object{
		"json": `{"a":1,"b":2,"c":3,"d":4,"e":5}`,
		"yaml": `id: 1
uuid: 3d877494-e7d4-48e3-aa7a-164373a7920d
name: He3
age: 26
isActive: true
registered: 2020-02-03T06:00:03 -08:00
tags:
  - tools
  - development
language:
  - id: 0
    name: English
  - id: 1
    name: Español
  - id: 2
    name: Chinese
`,
	})
	fmt.Println(result)
}

func TestPythonMikan(t *testing.T) {
	bangumiCache := cache.NewBolt()
	bangumiCache.Open("../../../test/testdata/bolt_sub.bolt")
	db := cache.NewBolt()
	db.Open("data/test.db")
	anidata.Init(&anidata.Options{
		Cache:        db,
		BangumiCache: bangumiCache,
	})
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "/Users/wetor/GoProjects/AnimeGo/assets/plugin/filter/AnimeGoHelperParser/mikan_tool.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})
	result := p.Run("main", models.Object{
		"url": `https://mikanani.me/Home/Episode/1672d44040d63b380ba0018d9e76e0fa0db18906`,
	})
	fmt.Println(result)
}

func TestPythonConfig(t *testing.T) {
	lib.Init()
	os.Setenv("ANIMEGO_VERSION", "0.6.8")
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "testdata/config.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "test",
				SkipSchemaCheck: true,
			},
		},
	})

	result := p.Run("test", models.Object{
		"test": true,
	})
	fmt.Println(result)
}

func TestPythonMikanTool(t *testing.T) {
	lib.Init()
	os.Setenv("ANIMEGO_VERSION", "0.6.8")

	db := cache.NewBolt()
	db.Open("data/bolt.db")
	bangumiCache := cache.NewBolt()
	bangumiCache.Open("../../../test/testdata/bolt_sub.bolt")
	anidata.Init(&anidata.Options{
		Cache:        db,
		BangumiCache: bangumiCache,
	})

	rss := mikanRss.NewRss("", "")
	items := rss.Parse("testdata/Mikan.xml")
	fmt.Println(len(items))
	fmt.Println("===========")

	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "/Users/wetor/GoProjects/AnimeGo/assets/plugin/filter/mikan_tool.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})

	result := p.Run("main", models.Object{
		"feedItems": items,
	})
	fmt.Println(result)
}

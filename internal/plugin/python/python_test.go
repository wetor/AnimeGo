package python_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	constant.PluginPath = "testdata"
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	gpython.Init()
	m.Run()
	fmt.Println("end")
}

func TestLib_log(t *testing.T) {
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "test_log.py",
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

func TestPythonFunction(t *testing.T) {
	lib.Init()
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "test.py",
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
		File: "var.py",
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
		File: "json.py",
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
name: wetor
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

func TestPythonConfig(t *testing.T) {
	lib.Init()
	os.Setenv("ANIMEGO_VERSION", "0.6.8")
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "config.py",
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

func TestPythonParseMikan(t *testing.T) {
	lib.Init()

	db := cache.NewBolt()
	db.Open("data/bolt.db")
	bangumiCache := cache.NewBolt()
	bangumiCache.Open("../../../test/testdata/bolt_sub.bolt")
	anidata.Init(&anidata.Options{
		Cache:        db,
		BangumiCache: bangumiCache,
	})

	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "mikan.py",
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})

	p.Run("main", nil)
}

func TestPythonMikanTool(t *testing.T) {
	savePluginPath := constant.PluginPath
	constant.PluginPath = "../../../assets/plugin"
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
	rss := mikanRss.NewRss(&mikanRss.Options{File: "testdata/Mikan.xml"})
	items := rss.Parse()
	fmt.Println(len(items))
	fmt.Println("===========")

	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: "filter/mikan_tool.py",
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
	fmt.Println(result["data"])
	constant.PluginPath = savePluginPath
}

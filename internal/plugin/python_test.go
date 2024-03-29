package plugin_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "python"

var mikanInst *mikan.Mikan

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	db := cache.NewBolt()
	db.Open("data/bolt.db")
	mikanInst = mikan.NewMikan(&mikan.Options{
		Cache: db,
	})
	m.Run()

	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestLib_log(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})
	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "test_log.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         "main",
				ParamsSchema: []string{"title"},
				ResultSchema: []string{"result"},
			},
		},
	})
	result, _ := p.Run("main", map[string]any{
		"title": "【悠哈璃羽字幕社】 [明日同学的水手服_Akebi-chan no Sailor-fuku] [01-12] [x264 1080p][CHT]",
	})
	fmt.Println(result)
}

func TestPythonFunction(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})
	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "test.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
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
	result, _ := p.Run("main", map[string]any{
		"params": []int{1, 2, 3},
	})
	fmt.Println(result)

	p.Run("test", map[string]any{
		"test": true,
	})
}

func TestPythonVariable(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})
	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "var.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
		VarSchema: []*pkgPlugin.VarSchemaOptions{
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
	result, _ := p.Run("main", map[string]any{
		"params": []int{1, 2, 3},
	})
	fmt.Println(result)

}

func TestPythonJson(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})
	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "json.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})
	result, _ := p.Run("main", map[string]any{
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
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})
	os.Setenv("ANIMEGO_VERSION", "0.6.8")
	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "config.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:            "test",
				SkipSchemaCheck: true,
			},
		},
	})

	result, _ := p.Run("test", map[string]any{
		"test": true,
	})
	fmt.Println(result)
}

func TestPythonParseMikan(t *testing.T) {
	test.HookGetWriter(testdata, nil)
	defer test.UnHook()
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
		Mikan: mikanInst,
	})
	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "mikan.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})

	p.Run("main", nil)
}

func TestPythonMikanTool(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
		Mikan: mikanInst,
	})
	os.Setenv("ANIMEGO_VERSION", "0.6.8")

	rss := feed.NewRss()
	items, _ := rss.ParseFile("testdata/Mikan.xml")
	fmt.Println(len(items))
	fmt.Println("===========")

	p := &python.Python{}
	p.Load(&pkgPlugin.LoadOptions{
		File: "filter/mikan_tool.py",
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:            "filter_all",
				SkipSchemaCheck: true,
			},
		},
	})

	result, _ := p.Run("filter_all", map[string]any{
		"items": items,
	})
	fmt.Println(result["data"])
}

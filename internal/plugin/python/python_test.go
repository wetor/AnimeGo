package python_test

import (
	"fmt"
	"github.com/brahma-adshonor/gohook"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/pkg/request"
	"io"
	"os"
	"path"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
)

func HookGetWriter(uri string, w io.Writer) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := path.Base(uri)
	jsonData, err := os.ReadFile(path.Join("testdata", id+".html"))
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	_ = gohook.Hook(request.GetWriter, HookGetWriter, nil)

	db := cache.NewBolt()
	db.Open("data/bolt.db")
	bangumiCache := cache.NewBolt()
	bangumiCache.Open("../../../test/testdata/bolt_sub.bolt")
	anidata.Init(&anidata.Options{
		Cache:        db,
		BangumiCache: bangumiCache,
	})

	m.Run()

	db.Close()
	bangumiCache.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestLib_log(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "testdata",
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
	result := p.Run("main", map[string]any{
		"title": "【悠哈璃羽字幕社】 [明日同学的水手服_Akebi-chan no Sailor-fuku] [01-12] [x264 1080p][CHT]",
	})
	fmt.Println(result)
}

func TestPythonFunction(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "testdata",
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
	result := p.Run("main", map[string]any{
		"params": []int{1, 2, 3},
	})
	fmt.Println(result)

	p.Run("test", map[string]any{
		"test": true,
	})
}

func TestPythonVariable(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "testdata",
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
	result := p.Run("main", map[string]any{
		"params": []int{1, 2, 3},
	})
	fmt.Println(result)

}

func TestPythonJson(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "testdata",
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
	result := p.Run("main", map[string]any{
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
		Path: "testdata",
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

	result := p.Run("test", map[string]any{
		"test": true,
	})
	fmt.Println(result)
}

func TestPythonParseMikan(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "testdata",
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
		Path: "../../../assets/plugin",
	})
	os.Setenv("ANIMEGO_VERSION", "0.6.8")

	rss := mikanRss.NewRss(&mikanRss.Options{File: "testdata/Mikan.xml"})
	items := rss.Parse()
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

	result := p.Run("filter_all", map[string]any{
		"items": items,
	})
	fmt.Println(result["data"])
}

package plugin_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	filterPlugin "github.com/wetor/AnimeGo/internal/animego/filter/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "filter"

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})
	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestPython_Filter(t *testing.T) {
	rss := mikanRss.NewRss(&mikanRss.Options{File: test.GetDataPath(testdata, "Mikan.xml")})
	items, _ := rss.Parse()
	fmt.Println(len(items))
	p := filterPlugin.NewFilterPlugin(&models.Plugin{
		Enable: true,
		Type:   "py",
		File:   "filter.py",
	})
	result, _ := p.FilterAll(items)
	assert.Equal(t, 4, len(result))
	for _, r := range result {
		fmt.Println(r)
	}
}

func TestPython_Filter2(t *testing.T) {
	list := []*models.FeedItem{
		{
			Name: "0000",
		},
		{
			Name: "1108011",
		},
		{
			Name: "2222",
		},
		{
			Name: "3333",
		},
	}
	p := filterPlugin.NewFilterPlugin(&models.Plugin{
		Enable: true,
		Type:   "py",
		File:   "test_re.py",
	})
	result, _ := p.FilterAll(list)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "1108011", result[0].Name)
}

func TestPython_Filter3(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	rss := mikanRss.NewRss(&mikanRss.Options{File: test.GetDataPath(testdata, "Mikan.xml")})
	items, _ := rss.Parse()
	fmt.Println(len(items))
	fmt.Println("===========")
	p := filterPlugin.NewFilterPlugin(&models.Plugin{
		Enable: true,
		Type:   "py",
		File:   "filter/pydemo.py",
	})
	result, _ := p.FilterAll(items)
	assert.Equal(t, 9, len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}
}

func TestPython_Filter4(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	list := []*models.FeedItem{
		{
			Name: "0000",
		},
		{
			Name: "1108011",
		},
		{
			Name: "2222",
		},
		{
			Name: "[梦蓝字幕组]New Doraemon 哆啦A梦新番[716][2022.07.23][AVC][1080P][GB_JP]",
		},
	}
	p := filterPlugin.NewFilterPlugin(&models.Plugin{
		Enable: true,
		Type:   "py",
		File:   "filter/default.py",
	})
	result, _ := p.FilterAll(list)
	assert.Equal(t, 4, len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}
}

func TestPython_Filter5(t *testing.T) {
	db := cache.NewBolt()
	db.Open("data/bolt.db")
	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open(test.GetDataPath("", "bolt_sub.bolt"))
	anidata.Init(&anidata.Options{
		Cache:        db,
		BangumiCache: bangumiCache,
	})

	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	rss := mikanRss.NewRss(&mikanRss.Options{File: test.GetDataPath(testdata, "Mikan.xml")})
	items, _ := rss.Parse()
	fmt.Println(len(items))
	fmt.Println("===========")
	p := filterPlugin.NewFilterPlugin(&models.Plugin{
		Enable: true,
		Type:   "py",
		File:   "filter/mikan_tool.py",
	})
	result, _ := p.FilterAll(items)
	assert.Equal(t, 13, len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}
	bangumiCache.Close()
	db.Close()
}

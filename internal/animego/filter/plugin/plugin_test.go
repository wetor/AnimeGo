package plugin

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/feed"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/internal/utils"
)

func TestJavaScript_Filter(t *testing.T) {
	_ = utils.CreateMutiDir("data")
	feed.Init(&feed.Options{
		TempPath: "data",
	})
	rss := mikanRss.NewRss("https://mikanani.me/RSS/MyBangumi?token=ky5DTt%2fMyAjCH2oKEN81FQ%3d%3d", "Mikan")
	items := rss.Parse()
	fmt.Println(len(items))
	js := NewPluginFilter(&javascript.JavaScript{}, []string{"testdata/test.js"})
	result := js.Filter(items)
	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}

}

func TestJavaScript_Filter2(t *testing.T) {
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
	js := NewPluginFilter(&javascript.JavaScript{}, []string{
		"testdata/regexp.js",
	})
	result := js.Filter(list)
	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}
}

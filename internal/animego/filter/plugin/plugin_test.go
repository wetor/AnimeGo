package plugin

import (
	"fmt"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestJavaScript_Filter(t *testing.T) {
	test.TestInit()

	feed := mikanRss.NewRss(store.Config.Setting.Feed.Mikan.Url, store.Config.Setting.Feed.Mikan.Name)
	items := feed.Parse()
	fmt.Println(len(items))
	js := NewPluginFilter(&javascript.JavaScript{}, store.Config.Filter.JavaScript)
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
		"/Users/wetor/GoProjects/AnimeGo/data/plugin/filter/regexp.js",
	})
	result := js.Filter(list)
	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}
}

package javascript

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/mikan"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	pkgAnisource "github.com/wetor/AnimeGo/pkg/anisource"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestJavaScript_Filter(t *testing.T) {
	test.TestInit()

	anisource.Init(&pkgAnisource.Options{
		Cache: store.Cache,
	})

	feed := mikanRss.NewRss(store.Config.Setting.Feed.Mikan.Url, store.Config.Setting.Feed.Mikan.Name)
	items, _ := feed.Parse()
	fmt.Println(len(items))
	js := &JavaScript{}
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
	js := &JavaScript{
		ScriptFile: []string{
			"/Users/wetor/GoProjects/AnimeGo/data/plugin/filter/regexp.js",
		},
	}
	result := js.Filter(list)
	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}
}

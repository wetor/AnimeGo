package javascript

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/mikan"
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

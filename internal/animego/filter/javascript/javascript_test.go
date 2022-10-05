package javascript

import (
	"AnimeGo/internal/animego/anisource"
	mikanRss "AnimeGo/internal/animego/feed/mikan"
	"AnimeGo/internal/store"
	pkgAnisource "AnimeGo/pkg/anisource"
	"AnimeGo/test"
	"fmt"
	"testing"
)

func TestJavaScript_Filter(t *testing.T) {
	test.TestInit()

	anisource.Init(&pkgAnisource.Options{
		Cache:   store.Cache,
		Proxy:   store.Config.Proxy(),
		Timeout: store.Config.HttpTimeoutSecond,
		Retry:   store.Config.HttpRetryNum,
	})

	feed := mikanRss.NewRss(store.Config.RssMikan().Url, store.Config.RssMikan().Name)
	items := feed.Parse()
	fmt.Println(len(items))
	js := &JavaScript{}
	result := js.Filter(items)
	fmt.Println(len(result))
	for _, r := range result {
		fmt.Println(r.Name)
	}

}

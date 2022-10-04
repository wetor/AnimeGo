package javascript

import (
	"AnimeGo/internal/animego/anisource"
	mikanRss "AnimeGo/internal/animego/feed/mikan"
	"AnimeGo/internal/logger"
	"AnimeGo/internal/store"
	pkgAnisource "AnimeGo/pkg/anisource"
	"fmt"
	"testing"
)

func TestJavaScript_Filter(t *testing.T) {
	logger.Init()
	defer logger.Flush()
	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/AnimeGo/data/config/conf.yaml",
	})

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

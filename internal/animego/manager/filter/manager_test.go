package filter

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/animego/anisource/mikan"
	mikanRss "AnimeGo/internal/animego/feed/mikan"
	"AnimeGo/internal/animego/filter"
	"AnimeGo/internal/store"
	pkgAnisource "AnimeGo/pkg/anisource"
	"AnimeGo/test"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	test.TestInit()

	m.Run()
	fmt.Println("end")
}

func TestManager_UpdateFeed(t *testing.T) {
	anisource.Init(&pkgAnisource.Options{
		Cache:   store.Cache,
		Proxy:   store.Config.Proxy(),
		Timeout: store.Config.HttpTimeoutSecond,
		Retry:   store.Config.HttpRetryNum,
	})

	rss := mikanRss.NewRss(store.Config.RssMikan().Url, store.Config.RssMikan().Name)
	mk := mikan.MikanAdapter{ThemoviedbKey: store.Config.KeyTmdb()}
	m := NewManager(&filter.Default{}, rss, mk, nil)

	exit := make(chan bool)

	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)

	go func() {
		time.Sleep(30 * time.Second)
		cancel()
		exit <- true
	}()

	//time.Sleep(120 * time.Second)

	<-exit
}

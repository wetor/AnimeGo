package filter

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/store"
	pkgAnisource "github.com/wetor/AnimeGo/pkg/anisource"
	"github.com/wetor/AnimeGo/test"
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
		Cache: store.Cache,
	})

	rss := mikanRss.NewRss(store.Config.Setting.Feed.Mikan.Url, store.Config.Setting.Feed.Mikan.Name)
	mk := mikan.MikanAdapter{ThemoviedbKey: store.Config.Setting.Key.Themoviedb}
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

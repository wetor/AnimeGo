package filter_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/animego/manager"
	filterMgr "github.com/wetor/AnimeGo/internal/animego/manager/filter"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
)

const ThemoviedbKey = "d3d8430aefee6c19520d0f7da145daf5"

var (
	wg sync.WaitGroup

	ctx, cancel = context.WithCancel(context.Background())
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	b := cache.NewBolt()
	b.Open("data/test.db")

	anisource.Init(&anisource.Options{
		Options: &anidata.Options{
			Cache: b,
		},
	})

	manager.Init(&manager.Options{
		Filter: manager.Filter{
			MultiGoroutineMax:     0,
			MultiGoroutineEnabled: false,
			UpdateDelayMinute:     10,
			DelaySecond:           2,
		},
		WG: &wg,
	})

	m.Run()
	fmt.Println("end")
}

func TestManager_UpdateFeed(t *testing.T) {

	rss := mikanRss.NewRss("https://mikanani.me/RSS/MyBangumi?token=ky5DTt%2fMyAjCH2oKEN81FQ%3d%3d", "Mikan")
	mk := mikan.Mikan{ThemoviedbKey: ThemoviedbKey}
	m := filterMgr.NewManager(&filter.Default{}, rss, mk, nil)

	m.Start(ctx)

	go func() {
		time.Sleep(10 * time.Second)
		cancel()
		wg.Done()
	}()

	wg.Wait()
}

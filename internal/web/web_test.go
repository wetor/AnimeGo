package web_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/configs"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/web"
	webapi "github.com/wetor/AnimeGo/internal/web/api"
	"github.com/wetor/AnimeGo/internal/web/websocket"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

type MockFilter struct{}

func (m *MockFilter) Update(ctx context.Context, items []*models.FeedItem, b, c bool) error {
	return nil
}

type MockDatabase struct{}

func (m *MockDatabase) Delete(data any) error { return nil }

var (
	ctx = context.Background()
	wg  = sync.WaitGroup{}
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")

	out, notify := logger.NewLogNotify()
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
		Out:   out,
	})

	bolt := cache.NewBolt()
	bolt.Open("data/bolt.db")
	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open(test.GetDataPath("", "bolt_sub.bolt"))
	bangumiCache.Add("bangumi_sub")
	BangumiCacheMutex := sync.Mutex{}

	config := configs.DefaultConfig()

	web.Init(&web.Options{
		ApiOptions: &webapi.Options{
			Ctx:                  ctx,
			AccessKey:            "animego123",
			Cache:                bolt,
			Config:               config,
			BangumiCache:         bangumiCache,
			BangumiCacheLock:     &BangumiCacheMutex,
			FilterManager:        &MockFilter{},
			DatabaseCacheDeleter: &MockDatabase{},
		},
		WebSocketOptions: &websocket.Options{
			Notify: notify,
			WG:     &wg,
		},
		Host:  config.WebApi.Host,
		Port:  config.WebApi.Port,
		WG:    &wg,
		Debug: true,
	})

	m.Run()

	bolt.Close()
	bangumiCache.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestRun(t *testing.T) {
	t.SkipNow()
	web.Run(ctx)

	go func() {
		i := 0
		for i < 1000 {
			time.Sleep(1 * time.Second)
			if logger.GetLogNotify() >= logger.NotifyEnabled {
				log.Debugf("日志输出：%d", i)
				i++
			}

		}
	}()
	wg.Wait()
}

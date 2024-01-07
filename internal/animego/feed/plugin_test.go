package feed_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/assets"
	feedPlugin "github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

var s *schedule.Schedule

type MockFilterManager struct {
}

func (m *MockFilterManager) Update(ctx context.Context, items []*models.FeedItem, b, c bool) error {
	for _, item := range items {
		fmt.Println(item)
	}
	return nil
}

func BeforePlugin() {
	fmt.Println("begin")
	constant.CachePath = "data"
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	test.HookGetString(testdata, func(uri string) string {
		u, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		bgmId := u.Query().Get("bangumiId")
		subId := u.Query().Get("subgroupid")
		return bgmId + "_" + subId + ".xml"
	})

	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
		Feed:  feedPlugin.NewRss(),
	})

	gpython.Init()
	_ = utils.CreateMutiDir("data")

	wg := sync.WaitGroup{}
	s = schedule.NewSchedule(&schedule.Options{
		WG: &wg,
	})
}

func After() {
	test.UnHook()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestNewSchedule3_feed(t *testing.T) {
	BeforePlugin()
	defer After()

	_ = feedPlugin.AddFeedTasks(s, []models.Plugin{
		{
			Enable: true,
			Type:   "builtin",
			File:   "builtin_mikan_rss.py",
			Vars: map[string]any{
				"__url__":  "https://mikanani.me/RSS/Bangumi?bangumiId=2822&subgroupid=370",
				"__cron__": "0/3 * * * * ?",
			},
		},
	}, &MockFilterManager{}, context.Background())

	s.Start(context.Background())
	time.Sleep(3 * time.Second)
	s.Delete("test")
}

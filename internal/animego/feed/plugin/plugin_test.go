package plugin_test

import (
	"context"
	"fmt"
	"github.com/brahma-adshonor/gohook"
	feedPlugin "github.com/wetor/AnimeGo/internal/animego/feed/plugin"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/third_party/gpython"
	"net/url"
	"os"
	"path"
	"sync"
	"testing"
	"time"
)

var s *schedule.Schedule

type MockFilterManager struct {
}

func (m *MockFilterManager) Update(ctx context.Context, items []*models.FeedItem) {
	for _, item := range items {
		fmt.Println(item)
	}
}

func GetString(uri string, args ...interface{}) (string, error) {
	log.Infof("Mock HTTP GET %s, header %s", uri, args)
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	bgm_id := u.Query().Get("bangumiId")
	sub_id := u.Query().Get("subgroupid")
	jsonData, err := os.ReadFile(path.Join("testdata", bgm_id+"_"+sub_id+".xml"))
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	constant.CachePath = "data"
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	_ = gohook.Hook(request.GetString, GetString, nil)

	plugin.Init(&plugin.Options{
		Path: "../../assets/plugin",
	})

	gpython.Init()
	lib.Init()
	_ = utils.CreateMutiDir("data")

	wg := sync.WaitGroup{}
	s = schedule.NewSchedule(&schedule.Options{
		WG: &wg,
	})
	m.Run()
	wg.Done()
	fmt.Println("end")
}

func TestNewSchedule3_feed(t *testing.T) {
	feedPlugin.AddFeedTasks(s, []models.Plugin{
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

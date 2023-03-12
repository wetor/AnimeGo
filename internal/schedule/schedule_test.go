package schedule_test

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
	"github.com/wetor/AnimeGo/pkg/utils"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

var s *schedule.Schedule

func TestMain(m *testing.M) {
	fmt.Println("begin")
	constant.CachePath = "data"
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
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

func TestNewSchedule(t *testing.T) {
	b := cache.NewBolt()
	b.Open("data/bolt_sub.db")
	mutex := sync.Mutex{}

	s.Add(&schedule.AddTaskOptions{
		Name:     "bangumi",
		StartRun: true,
		Task: task.NewBangumiTask(&task.BangumiOptions{
			Cache:      b,
			CacheMutex: &mutex,
		}),
	})
	s.Start(context.Background())
	time.Sleep(1 * time.Second)
	s.Delete("bangumi")
}

func TestNewSchedule2(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "../../assets/plugin",
	})
	s.Add(&schedule.AddTaskOptions{
		Name:     "test",
		StartRun: true,
		Task: task.NewScheduleTask(&task.ScheduleOptions{
			Plugin: &models.Plugin{
				Enable: true,
				Type:   python.Type,
				File:   "schedule/refresh.py",
				Vars: models.Object{
					"__name__": "Vars_Test",
				},
				Args: models.Object{
					"Args_Test": 13213,
				},
			},
		}),
		Args: models.Object{
			"Args_Test": "测试",
		},
		Vars: models.Object{
			"__name__": "outer_Vars_Test",
		},
	})
	s.Start(context.Background())
	time.Sleep(11 * time.Second)
	s.Delete("test")
}

func TestNewSchedule3_feed(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path: "../../assets/plugin",
	})
	s.Add(&schedule.AddTaskOptions{
		Name:     "test",
		StartRun: true,
		Task: task.NewFeedTask(&task.FeedOptions{
			Plugin: &models.Plugin{
				Enable: true,
				Type:   python.Type,
				File:   "feed/mikan_rss.py",
			},
			Callback: func(items []*models.FeedItem) {
				for i, item := range items {
					fmt.Println(i, "download: ", item)
				}
			},
		}),
	})
	s.Start(context.Background())
	time.Sleep(100 * time.Second)
	s.Delete("test")
}

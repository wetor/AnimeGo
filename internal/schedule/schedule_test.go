package schedule_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/plugin/lib"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

var (
	s  *schedule.Schedule
	wg = sync.WaitGroup{}
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	constant.CachePath = "data"
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})

	gpython.Init()
	lib.Init(&lib.Options{
		Feed: feed.NewRss(),
	})
	_ = utils.CreateMutiDir("data")

	s = schedule.NewSchedule(&schedule.Options{
		WG: &wg,
	})
	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestNewSchedule(t *testing.T) {
	t.Skip("跳过BangumiCache更新测试")
	b := cache.NewBolt()
	b.Open("data/bolt_sub.db")
	mutex := sync.Mutex{}

	s.Add(&schedule.AddTaskOptions{
		Name:     "bangumi",
		StartRun: true,
		Task: schedule.NewBangumiTask(&schedule.BangumiOptions{
			Cache:      b,
			CacheMutex: &mutex,
		}),
	})
	s.Start(context.Background())
	time.Sleep(1 * time.Second)
	s.Delete("bangumi")
	wg.Done()
	b.Close()
}

func TestNewSchedule2(t *testing.T) {
	tt, _ := schedule.NewScheduleTask(&schedule.PluginOptions{
		Plugin: &models.Plugin{
			Enable: true,
			Type:   "python",
			File:   "schedule/refresh.py",
			Vars: models.Object{
				"__name__": "Vars_Test",
			},
			Args: models.Object{
				"Args_Test": 13213,
			},
		},
	})
	s.Add(&schedule.AddTaskOptions{
		Name:     "test",
		StartRun: false,
		Task:     tt,
		Args: models.Object{
			"Args_Test1": "测试",
		},
		Vars: models.Object{
			"name":     "outer_Vars_Test",
			"__cron__": "*/3 * * * * ?",
		},
	})
	s.Start(context.Background())
	time.Sleep(7 * time.Second)
	s.Delete("test")
	wg.Done()
}

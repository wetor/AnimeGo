package schedule_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
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
		Path: "testdata",
	})

	gpython.Init()
	lib.Init()
	_ = utils.CreateMutiDir("data")

	s = schedule.NewSchedule(&schedule.Options{
		WG: &wg,
	})
	m.Run()
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
		Task: task.NewBangumiTask(&task.BangumiOptions{
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
	s.Add(&schedule.AddTaskOptions{
		Name:     "test",
		StartRun: false,
		Task: task.NewScheduleTask(&task.ScheduleOptions{
			Plugin: &models.Plugin{
				Enable: true,
				Type:   "python",
				File:   "refresh.py",
				Vars: models.Object{
					"__name__": "Vars_Test",
				},
				Args: models.Object{
					"Args_Test": 13213,
				},
			},
		}),
		Args: models.Object{
			"Args_Test1": "测试",
		},
		Vars: models.Object{
			"__name__": "outer_Vars_Test",
			"__cron__": "*/3 * * * * ?",
		},
	})
	s.Start(context.Background())
	time.Sleep(7 * time.Second)
	s.Delete("test")
	wg.Done()
}

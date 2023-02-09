package schedule_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

var s *schedule.Schedule

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "log/log.log",
		Debug: true,
	})
	gpython.Init()
	lib.InitLog()
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
			DBPath:     "data",
			Cache:      b,
			CacheMutex: &mutex,
		}),
	})
	s.Start(context.Background())
	time.Sleep(1 * time.Second)
	s.Delete("bangumi")
}

func TestNewSchedule2(t *testing.T) {
	constant.PluginPath = "/Users/wetor/GoProjects/AnimeGo/assets/plugin"
	s.Add(&schedule.AddTaskOptions{
		Name:     "test",
		StartRun: true,
		Task: task.NewPluginTask(&task.PluginOptions{
			Plugin: &models.Plugin{
				Enable: true,
				Type:   python.Type,
				File:   "schedule/refresh.py",
			},
		}),
		Params: []interface{}{models.Object{"aaa": "测试内容"}},
	})
	s.Start(context.Background())
	time.Sleep(11 * time.Second)
	s.Delete("test")
}

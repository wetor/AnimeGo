package schedule_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/third_party/gpython"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/plugin/python"

	"github.com/wetor/AnimeGo/internal/models"

	"github.com/wetor/AnimeGo/pkg/log"

	"github.com/wetor/AnimeGo/internal/utils"

	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/cache"
)

func TestNewSchedule(t *testing.T) {
	log.Init(&log.Options{
		File:  "log/log.log",
		Debug: true,
	})
	constant.PluginPath = "/Users/wetor/GoProjects/AnimeGo/assets/plugin"

	gpython.Init()
	lib.InitLog()

	_ = utils.CreateMutiDir("task/data")

	b := cache.NewBolt()
	b.Open("task/data/bolt_sub.db")
	mutex := sync.Mutex{}
	wg := sync.WaitGroup{}
	schedule.Init(&schedule.Options{
		Options: &task.Options{
			DBDir:            "task/data",
			BangumiCache:     b,
			BangumiCacheLock: &mutex,
		},
		WG: &wg,
	})
	s := schedule.NewSchedule()
	s.Add("test", task.NewPluginTask(&task.PluginOptions{
		Cron: "*/5 * * * * ?",
		Plugin: &models.Plugin{
			Enable: true,
			Type:   python.Type,
			File:   "schedule/refresh.py",
		},
	}))
	for _, ts := range s.List() {
		fmt.Println(ts)
	}
	s.Start(context.Background())
	time.Sleep(11 * time.Second)
	s.Delete("test")
	time.Sleep(1 * time.Second)
}

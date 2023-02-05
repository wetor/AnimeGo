package schedule_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/cache"
)

func TestNewSchedule(t *testing.T) {
	wg := sync.WaitGroup{}
	b := cache.NewBolt()
	b.Open("task/data/bolt_sub.db")
	mutex := sync.Mutex{}

	schedule.Init(&schedule.Options{
		Options: &task.Options{
			DBDir:            "task/data",
			BangumiCache:     b,
			BangumiCacheLock: &mutex,
		},
		WG: &wg,
	})
	s := schedule.NewSchedule()
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	s.Add("test", task.NewJSPluginTask(&parser))
	for _, ts := range s.List() {
		fmt.Println(ts)
	}
	s.Start(context.Background())
	time.Sleep(11 * time.Second)
	s.Delete("test")
	time.Sleep(5 * time.Second)

}

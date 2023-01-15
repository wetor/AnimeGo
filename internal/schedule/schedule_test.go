package schedule

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/cache"
)

func TestNewSchedule(t *testing.T) {
	wg := sync.WaitGroup{}
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	b := cache.NewBolt()
	b.Open("task/data/bolt_sub.db")
	mutex := sync.Mutex{}

	Init(&Options{
		Options: &task.Options{
			DBDir:            "task/data",
			BangumiCache:     b,
			BangumiCacheLock: &mutex,
		},
		WG: &wg,
	})
	s := NewSchedule()
	s.Add("test", task.NewJSPluginTask(&s.parser))
	for _, ts := range s.List() {
		fmt.Println(ts)
	}
	s.Start(context.Background())
	time.Sleep(11 * time.Second)
	s.Delete("test")
	time.Sleep(5 * time.Second)

}

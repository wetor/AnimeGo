package task

import (
	"fmt"
	"sync"
	"testing"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/pkg/cache"
)

func TestTask_Bangumi_Start(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	b := cache.NewBolt()
	b.Open("data/bolt_sub.db")
	mutex := sync.Mutex{}
	Init(&Options{
		DBDir:            "data",
		BangumiCache:     b,
		BangumiCacheLock: &mutex,
	})

	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	task := NewBangumiTask(&p)
	fmt.Println(task.NextTime())
	task.Run(true)

}

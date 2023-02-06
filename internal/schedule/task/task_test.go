package task

import (
	"fmt"
	"sync"
	"testing"

	"github.com/wetor/AnimeGo/pkg/cache"
)

func TestTask_Bangumi_Start(t *testing.T) {
	b := cache.NewBolt()
	b.Open("data/bolt_sub.db")
	mutex := sync.Mutex{}
	Init(&Options{
		DBDir:            "data",
		BangumiCache:     b,
		BangumiCacheLock: &mutex,
	})

	task := NewBangumiTask()
	fmt.Println(task.NextTime())
	task.Run(true)

}

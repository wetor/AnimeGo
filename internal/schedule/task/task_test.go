package task

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestBangumiTask_Start(t *testing.T) {
	test.TestInit()
	dir := "/Users/wetor/GoProjects/AnimeGo/data/cache"
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	task := NewBangumiTask(dir, &p)
	fmt.Println(task.NextTime())
	task.Run(true)
}

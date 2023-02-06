package schedule

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

const (
	RetryNum  = 3 // 失败重试次数
	RetryWait = 1 // 失败重试等待时间，秒
)

var (
	WG *sync.WaitGroup
)

type Options struct {
	*task.Options
	WG *sync.WaitGroup
}

func Init(opts *Options) {
	WG = opts.WG
	task.Init(opts.Options)
}

type Schedule struct {
	tasks   map[string]api.Task
	task2id map[string]cron.EntryID
	crontab *cron.Cron
}

func NewSchedule() *Schedule {
	schedule := &Schedule{
		tasks:   make(map[string]api.Task),
		task2id: make(map[string]cron.EntryID),
	}
	schedule.crontab = cron.New(cron.WithSeconds(), cron.WithLogger(logger.NewCronLoggerAdapter()))

	schedule.Add("bangumi", task.NewBangumiTask())
	schedule.tasks["bangumi"].Run(true)
	return schedule
}

func (s *Schedule) Add(name string, task api.Task) {
	id, err := s.crontab.AddFunc(task.Cron(), func() {
		log.Infof("[定时任务] %s 开始执行", task.Name())
		success := false
		for i := 0; i < RetryNum; i++ {
			try.This(func() {
				task.Run(false)
				success = true
			}).Catch(func(err try.E) {
				log.Debugf("", err)
				if i == RetryNum-1 {
					log.Warnf("[定时任务] %s 第%d次执行失败", task.Name(), i+1)
				} else {
					log.Warnf("[定时任务] %s 第%d次执行失败，%d 秒后重新执行", task.Name(), i+1, RetryWait)
				}
				utils.Sleep(RetryWait, context.Background())
			})
			if success {
				log.Infof("[定时任务] %s 执行完毕，下次执行时间: %s", task.Name(), task.NextTime())
				break
			}
		}
	})
	if err != nil {
		errors.NewAniErrorD(err).TryPanic()
	}
	s.tasks[name] = task
	s.task2id[name] = id
}

func (s *Schedule) Delete(name string) {
	s.crontab.Remove(s.task2id[name])
	delete(s.tasks, name)
	delete(s.task2id, name)
}

func (s *Schedule) List() []*task.TaskInfo {
	list := make([]*task.TaskInfo, 0, len(s.tasks))
	for name, task_ := range s.tasks {
		list = append(list, &task.TaskInfo{
			Name:  name,
			RunAt: task_.NextTime(),
			Cron:  task_.Cron(),
		})
	}
	return list
}

func (s *Schedule) Start(ctx context.Context) {
	s.crontab.Start()
	WG.Add(1)
	go func() {
		defer WG.Done()
		for {
			select {
			case <-ctx.Done():
				s.crontab.Stop()
				return
			}
		}
	}()
}

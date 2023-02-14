package schedule

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/logger"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

const (
	RetryNum  = 3 // 失败重试次数
	RetryWait = 1 // 失败重试等待时间，秒
)

type TaskInfo struct {
	Name    string
	Task    api.Task
	Params  []interface{}
	TaskRun func(params ...interface{})
}

type Schedule struct {
	wg      *sync.WaitGroup
	tasks   map[string]*TaskInfo
	task2id map[string]cron.EntryID
	crontab *cron.Cron
}

type AddTaskOptions struct {
	Name     string
	StartRun bool
	Params   []interface{}
	Task     api.Task
}

type Options struct {
	WG *sync.WaitGroup
}

func NewSchedule(opts *Options) *Schedule {
	schedule := &Schedule{
		wg:      opts.WG,
		tasks:   make(map[string]*TaskInfo),
		task2id: make(map[string]cron.EntryID),
	}
	schedule.crontab = cron.New(cron.WithSeconds(), cron.WithLogger(logger.NewCronLoggerAdapter()))
	return schedule
}

func (s *Schedule) Add(opts *AddTaskOptions) {
	t := &TaskInfo{
		Name:   opts.Name,
		Task:   opts.Task,
		Params: opts.Params,
		TaskRun: func(params ...interface{}) {
			log.Infof("[定时任务] %s 开始执行", opts.Task.Name())
			success := false
			for i := 0; i < RetryNum; i++ {
				try.This(func() {
					opts.Task.Run(params...)
					success = true
				}).Catch(func(err try.E) {
					log.Debugf("", err)
					if i == RetryNum-1 {
						log.Warnf("[定时任务] %s 第%d次执行失败", opts.Task.Name(), i+1)
					} else {
						log.Warnf("[定时任务] %s 第%d次执行失败，%d 秒后重新执行", opts.Task.Name(), i+1, RetryWait)
					}
					utils.Sleep(RetryWait, context.Background())
				})
				if success {
					log.Infof("[定时任务] %s 执行完毕，下次执行时间: %s", opts.Task.Name(), opts.Task.NextTime())
					break
				}
			}
		},
	}

	id, err := s.crontab.AddFunc(t.Task.Cron(), func() {
		t.TaskRun(t.Params...)
	})
	if err != nil {
		errors.NewAniErrorD(err).TryPanic()
	}
	if opts.StartRun {
		t.TaskRun(t.Params...)
	}
	s.tasks[t.Name] = t
	s.task2id[t.Name] = id
}

func (s *Schedule) Get(name string) *TaskInfo {
	return s.tasks[name]
}

func (s *Schedule) Delete(name string) {
	s.crontab.Remove(s.task2id[name])
	delete(s.tasks, name)
	delete(s.task2id, name)
}

func (s *Schedule) List() []*TaskInfo {
	list := make([]*TaskInfo, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	return list
}

func (s *Schedule) Start(ctx context.Context) {
	s.crontab.Start()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-ctx.Done():
				s.crontab.Stop()
				return
			}
		}
	}()
}

package schedule

import (
	"context"
	"sync"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/errors"
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
	tasks   map[string]task.Task
	task2id map[string]cron.EntryID
	crontab *cron.Cron
	parser  cron.Parser
}

func NewSchedule() *Schedule {
	schedule := &Schedule{
		tasks:   make(map[string]task.Task),
		task2id: make(map[string]cron.EntryID),
		parser:  cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}
	schedule.crontab = cron.New(cron.WithParser(schedule.parser))

	schedule.Add("bangumi", task.NewBangumiTask(&schedule.parser))

	return schedule
}

func (s *Schedule) Add(name string, task task.Task) {
	task.Run(true)
	id, err := s.crontab.AddFunc(task.Cron(), func() {
		task.Run(false)
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

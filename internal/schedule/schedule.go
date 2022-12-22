package schedule

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/errors"
	"path"
)

type Schedule struct {
	tasks   map[string]Task
	task2id map[string]cron.EntryID
	crontab *cron.Cron
	parser  cron.Parser
}

func NewSchedule() *Schedule {
	schedule := &Schedule{
		tasks:   make(map[string]Task),
		task2id: make(map[string]cron.EntryID),
		parser:  cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}
	schedule.crontab = cron.New(cron.WithParser(schedule.parser))

	schedule.tasks["bangumi"] = task.NewBangumiTask(path.Dir(store.Config.Advanced.Path.DbFile), &schedule.parser)
	// schedule.tasks["js_plugin"] = task.NewJSPluginTask(&schedule.parser)

	for name, task_ := range schedule.tasks {
		task_.Run(true)
		id, err := schedule.crontab.AddFunc(task_.Cron(), func() {
			task_.Run(false)
		})
		if err != nil {
			errors.NewAniErrorD(err).TryPanic()
		}
		schedule.task2id[name] = id
	}
	return schedule
}

func (s *Schedule) Add(name string, task Task) {
	s.tasks[name] = task
}

func (s *Schedule) Delete(name string) {
	delete(s.tasks, name)
	delete(s.task2id, name)
}

func (s *Schedule) List() []*TaskInfo {
	list := make([]*TaskInfo, 0, len(s.tasks))
	for name, task_ := range s.tasks {
		list = append(list, &TaskInfo{
			Name:  name,
			RunAt: task_.NextTime(),
			Cron:  task_.Cron(),
		})
	}
	return list
}

func (s *Schedule) Start(ctx context.Context) {
	s.crontab.Start()
	store.WG.Add(1)
	go func() {
		defer store.WG.Done()
		for {
			select {
			case <-ctx.Done():
				s.crontab.Stop()
				return
			}
		}
	}()
}

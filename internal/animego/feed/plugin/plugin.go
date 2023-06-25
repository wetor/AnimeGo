package plugin

import (
	"context"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

func AddFeedTasks(s *schedule.Schedule, plugins []models.Plugin, filterManager api.FilterManager, ctx context.Context) (err error) {
	for _, p := range plugins {
		if !p.Enable {
			continue
		}
		t, err := task.NewFeedTask(&task.FeedOptions{
			Plugin: &p,
			Callback: func(items []*models.FeedItem) error {
				return filterManager.Update(ctx, items, nil, false, false)
			},
		})
		if err != nil {
			return err
		}
		err = s.Add(&schedule.AddTaskOptions{
			Name:     xpath.Base(p.File),
			StartRun: false,
			Vars:     p.Vars,
			Args:     p.Args,
			Task:     t,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

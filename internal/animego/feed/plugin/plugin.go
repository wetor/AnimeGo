package plugin

import (
	"context"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const Builtin = "builtin"

func AddFeedTasks(s *schedule.Schedule, plugins []models.Plugin, filterManager api.FilterManager, ctx context.Context) {
	for _, p := range plugins {
		if !p.Enable {
			continue
		}
		s.Add(&schedule.AddTaskOptions{
			Name:     xpath.Base(p.File),
			StartRun: false,
			Vars:     p.Vars,
			Args:     p.Args,
			Task: task.NewFeedTask(&task.FeedOptions{
				Plugin: &p,
				Callback: func(items []*models.FeedItem) {
					filterManager.Update(ctx, items)
				},
			}),
		})
	}
}

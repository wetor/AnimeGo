package feed

import (
	"context"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"path"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/schedule"
)

func AddFeedTasks(s *schedule.Schedule, plugins []models.Plugin, filterManager api.FilterManager, ctx context.Context) (err error) {
	for _, p := range plugins {
		if !p.Enable {
			continue
		}
		t, err := schedule.NewFeedTask(&schedule.FeedOptions{
			Plugin: &p,
			Callback: func(items []*models.FeedItem) error {
				return filterManager.Update(ctx, items, false, false)
			},
		})
		if err != nil {
			return err
		}
		err = s.Add(&schedule.AddTaskOptions{
			Name:     path.Base(xpath.P(p.File)),
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

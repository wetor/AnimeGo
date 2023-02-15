package utils

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

func AddTasks(s *schedule.Schedule, plugins []models.Plugin) {
	for _, p := range plugins {
		s.Add(&schedule.AddTaskOptions{
			Name:     xpath.Base(p.File),
			StartRun: false,
			Task: task.NewPluginTask(&task.PluginOptions{
				Plugin: &p,
			}),
		})
	}
}

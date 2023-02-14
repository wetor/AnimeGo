package utils

import (
	"path/filepath"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/schedule/task"
)

func AddTasks(s *schedule.Schedule, plugins []models.Plugin) {
	for _, p := range plugins {
		s.Add(&schedule.AddTaskOptions{
			Name:     filepath.Base(p.File),
			StartRun: false,
			Task: task.NewPluginTask(&task.PluginOptions{
				Plugin: &p,
			}),
		})
	}
}

package api

import (
	"github.com/wetor/AnimeGo/internal/models"
	"time"
)

type Task interface {
	Cron() string
	NextTime() time.Time
	Name() string
	Run(args models.Object)
	SetVars(vars models.Object)
}

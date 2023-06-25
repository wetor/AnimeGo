package api

import (
	"time"

	"github.com/wetor/AnimeGo/internal/models"
)

type Task interface {
	Cron() string
	NextTime() time.Time
	Name() string
	Run(args models.Object) error
	SetVars(vars models.Object)
}

package schedule

import (
	"time"
)

type TaskInfo struct {
	Name  string
	RunAt time.Time
	Cron  string
}

type Task interface {
	Cron() string
	NextTime() time.Time
	Name() string
	Run(force bool)
}

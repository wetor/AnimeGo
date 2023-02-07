package api

import "time"

type Task interface {
	Cron() string
	NextTime() time.Time
	Name() string
	Run(params ...interface{})
}

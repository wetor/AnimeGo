package task

import (
	"github.com/robfig/cron/v3"
)

var SecondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

package schedule

import (
	"github.com/robfig/cron/v3"
)

const (
	VarName   = "__name__"
	VarCron   = "__cron__"
	VarUrl    = "__url__"
	VarHeader = "__header__"
	FuncRun   = "run"
	FuncParse = "parse"
)

var SecondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

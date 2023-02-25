package task

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/pkg/errors"
)

type ScheduleTask struct {
	parser *cron.Parser
	cron   string
	plugin api.Plugin
}

type ScheduleOptions struct {
	*models.Plugin
	Cron string
}

func NewScheduleTask(opts *ScheduleOptions) *ScheduleTask {
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: opts.File,
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            FuncRun,
				SkipSchemaCheck: true,
			},
		},
		Variables: []*models.PluginVariableOptions{
			{
				Name:     VarName,
				Nullable: true,
			},
			{
				Name: VarCron,
			},
		},
	})
	if len(opts.Cron) > 0 {
		_, err := SecondParser.Parse(opts.Cron)
		if err == nil {
			p.Set(VarCron, opts.Cron)
		}
	}
	return &ScheduleTask{
		parser: &SecondParser,
		cron:   p.Get(VarCron).(string),
		plugin: p,
	}
}

func (t *ScheduleTask) Name() string {
	name := t.plugin.Get(VarName)
	if name == nil {
		name = "Schedule"
	}
	return fmt.Sprintf("%v(%s-Plugin)", name, t.plugin.Type())
}

func (t *ScheduleTask) Cron() string {
	return t.plugin.Get(VarCron).(string)
}

func (t *ScheduleTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.cron)
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *ScheduleTask) Run(params ...interface{}) {
	t.plugin.Run(FuncRun, nil)
}

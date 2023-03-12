package task

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
)

type ScheduleTask struct {
	parser *cron.Parser
	cron   string
	plugin api.Plugin
	args   models.Object
}

type ScheduleOptions struct {
	*models.Plugin
}

func NewScheduleTask(opts *ScheduleOptions) *ScheduleTask {
	p := &python.Python{}
	p.Load(&plugin.LoadOptions{
		File: opts.File,
		Functions: []*plugin.FunctionOptions{
			{
				Name:            FuncRun,
				SkipSchemaCheck: true,
			},
		},
		Variables: []*plugin.VariableOptions{
			{
				Name:     VarName,
				Nullable: true,
			},
			{
				Name: VarCron,
			},
		},
	})
	for name, val := range opts.Plugin.Vars {
		if name == VarCron {
			_, err := SecondParser.Parse(val.(string))
			errors.NewAniErrorD(err).TryPanic()
		}
		p.Set(name, val)
	}
	return &ScheduleTask{
		parser: &SecondParser,
		cron:   p.Get(VarCron).(string),
		plugin: p,
		args:   opts.Args,
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

func (t *ScheduleTask) SetVars(vars models.Object) {
	for k, v := range vars {
		t.plugin.Set(k, v)
	}
}

func (t *ScheduleTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.cron)
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *ScheduleTask) Run(args models.Object) {
	for k, v := range t.args {
		if _, ok := args[k]; !ok {
			args[k] = v
		}
	}
	t.plugin.Run(FuncRun, args)
}

package task

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
)

type ScheduleTask struct {
	parser *cron.Parser
	plugin api.Plugin
	args   models.Object
}

type ScheduleOptions struct {
	*models.Plugin
}

func NewScheduleTask(opts *ScheduleOptions) *ScheduleTask {
	p := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    opts.Plugin,
		EntryFunc: FuncRun,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:            FuncRun,
				SkipSchemaCheck: true,
			},
		},
		VarSchema: []*pkgPlugin.VarSchemaOptions{
			{
				Name:     VarName,
				Nullable: true,
			},
			{
				Name: VarCron,
			},
		},
	})
	return &ScheduleTask{
		parser: &SecondParser,
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
	next, err := t.parser.Parse(t.Cron())
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

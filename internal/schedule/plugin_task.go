package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
)

type PluginTask struct {
	parser *cron.Parser
	plugin api.Plugin
	args   models.Object
}

type PluginOptions struct {
	*models.Plugin
}

func NewScheduleTask(opts *PluginOptions) (*PluginTask, error) {
	p, err := plugin.LoadPlugin(&plugin.LoadPluginOptions{
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
	if err != nil {
		return nil, err
	}
	return &PluginTask{
		parser: &SecondParser,
		plugin: p,
		args:   opts.Args,
	}, nil
}

func (t *PluginTask) Name() string {
	name, err := t.plugin.Get(VarName)
	if err != nil {
		log.Warnf("%s", err)
	}

	if name == nil {
		name = "Schedule"
	}
	return fmt.Sprintf("%v(%s-Plugin)", name, t.plugin.Type())
}

func (t *PluginTask) Cron() string {
	cronStr, err := t.plugin.Get(VarCron)
	if err != nil {
		log.Warnf("%s", err)
	}
	return cronStr.(string)
}

func (t *PluginTask) SetVars(vars models.Object) {
	for k, v := range vars {
		err := t.plugin.Set(k, v)
		if err != nil {
			log.Warnf("%s", err)
		}
	}
}

func (t *PluginTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.Cron())
	if err != nil {
		log.DebugErr(err)
	}
	return next.Next(time.Now())
}

func (t *PluginTask) Run(args models.Object) (err error) {
	for k, v := range t.args {
		if _, ok := args[k]; !ok {
			args[k] = v
		}
	}
	_, err = t.plugin.Run(FuncRun, args)
	if err != nil {
		return err
	}
	return nil
}

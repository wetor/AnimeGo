package task

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/wetor/AnimeGo/internal/constant"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

type PluginTask struct {
	parser *cron.Parser
	cron   string

	plugin api.Plugin
}

type PluginOptions struct {
	*models.Plugin
	Cron string
}

func NewPluginTask(opts *PluginOptions) *PluginTask {
	p := plugin.GetPlugin(opts.Type, plugin.Schedule)
	p.Load(&models.PluginLoadOptions{
		File: filepath.Join(constant.PluginPath, opts.File),
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "Run",
				SkipSchemaCheck: true,
			},
		},
		Variables: []*models.PluginVariableOptions{
			{
				Name:     "Name",
				Nullable: true,
			},
			{
				Name: "Cron",
			},
		},
	})
	if len(opts.Cron) > 0 {
		_, err := SecondParser.Parse(opts.Cron)
		if err == nil {
			p.Set("Cron", opts.Cron)
		}
	}
	return &PluginTask{
		parser: &SecondParser,
		cron:   p.Get("Cron").(string),
		plugin: p,
	}
}

func (t *PluginTask) Name() string {
	name := t.plugin.Get("Name")
	if name == nil {
		name = "NoName"
	}
	return fmt.Sprintf("%v(%s-Plugin)", t.plugin.Get("Name"), t.plugin.Type())
}

func (t *PluginTask) Cron() string {
	return t.plugin.Get("Cron").(string)
}

func (t *PluginTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.cron)
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *PluginTask) Run(params ...interface{}) {
	var obj models.Object
	var ok bool
	if len(params) >= 1 {
		if obj, ok = params[0].(models.Object); !ok {
			log.Debugf("[定时任务] %s-Plugin 参数错误: %v", t.plugin.Type(), params[0])
		}
	}
	t.plugin.Run("Run", obj)
}

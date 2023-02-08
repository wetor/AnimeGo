package task

import (
	"fmt"
	"path"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

type PluginTask struct {
	parser *cron.Parser
	cron   string

	file   string
	plugin api.Plugin
}

type PluginOptions struct {
	*models.Plugin
	Cron string
}

func NewPluginTask(opts *PluginOptions) *PluginTask {
	return &PluginTask{
		parser: &SecondParser,
		cron:   opts.Cron,
		file:   opts.File,
		plugin: plugin.GetPlugin(opts.Type, plugin.Schedule),
	}
}

func (t *PluginTask) Name() string {
	if t.plugin == nil {
		return "NoInit-Plugin"
	}
	return fmt.Sprintf("%s-Plugin", t.plugin.Type())
}

func (t *PluginTask) Cron() string {
	return t.cron
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
	t.plugin.Load(&models.PluginLoadOptions{
		File: path.Join(constant.PluginPath, t.file),
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})
	t.plugin.Run("main", obj)
}

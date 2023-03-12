package task

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/try"
)

type FeedTask struct {
	parser   *cron.Parser
	cron     string
	plugin   api.Plugin
	args     models.Object
	callback func([]*models.FeedItem)
}

type FeedOptions struct {
	*models.Plugin
	Callback func([]*models.FeedItem)
}

func NewFeedTask(opts *FeedOptions) *FeedTask {
	p := &python.Python{}
	p.Load(&models.PluginLoadOptions{
		File: opts.File,
		Functions: []*models.PluginFunctionOptions{
			{
				Name:         FuncParse,
				ParamsSchema: []string{"data"},
				ResultSchema: []string{"error", "items"},
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
			{
				Name: VarUrl,
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
	return &FeedTask{
		parser:   &SecondParser,
		cron:     p.Get(VarCron).(string),
		plugin:   p,
		args:     opts.Args,
		callback: opts.Callback,
	}
}

func (t *FeedTask) Name() string {
	name := t.plugin.Get(VarName)
	if name == nil {
		name = "Feed"
	}
	return fmt.Sprintf("%v(%s-Plugin)", name, t.plugin.Type())
}

func (t *FeedTask) Cron() string {
	return t.plugin.Get(VarCron).(string)
}

func (t *FeedTask) SetVars(vars models.Object) {
	for k, v := range vars {
		t.plugin.Set(k, v)
	}
}

func (t *FeedTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.cron)
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *FeedTask) Run(args models.Object) {
	url := t.plugin.Get(VarUrl).(string)
	data, err := request.GetString(url)
	if err != nil {
		log.Warnf("[Plugin] %s插件(%s)执行错误: 请求 %s 失败", t.plugin.Type(), FuncParse, url)
		log.Debugf("", err)
	}
	for k, v := range t.args {
		if _, ok := args[k]; !ok {
			args[k] = v
		}
	}
	args["data"] = data
	result := t.plugin.Run(FuncParse, args)
	if result["error"] != nil {
		log.Debugf("", errors.NewAniErrorD(result["error"]))
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", t.plugin.Type(), t.Name(), result["error"])
	}

	try.This(func() {
		itemsAny := result["items"].([]any)
		items := make([]*models.FeedItem, len(itemsAny))
		for i, item := range itemsAny {
			items[i] = &models.FeedItem{}
			utils.MapToStruct(item.(models.Object), items[i])
		}
		t.callback(items)
	}).Catch(func(err try.E) {
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", t.plugin.Type(), t.Name(), err)
		log.Debugf("", err)
	})

}

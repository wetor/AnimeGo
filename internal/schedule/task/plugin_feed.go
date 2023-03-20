package task

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/try"
	"github.com/wetor/AnimeGo/pkg/utils"
)

type FeedTask struct {
	parser   *cron.Parser
	plugin   api.Plugin
	args     models.Object
	callback func([]*models.FeedItem)
}

type FeedOptions struct {
	*models.Plugin
	Callback func([]*models.FeedItem)
}

func NewFeedTask(opts *FeedOptions) *FeedTask {
	p := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    opts.Plugin,
		EntryFunc: FuncParse,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         FuncParse,
				ParamsSchema: []string{"data"},
				ResultSchema: []string{"error", "items"},
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
			{
				Name: VarUrl,
			},
			{
				Name:     VarHeader,
				Nullable: true,
			},
		},
	})
	return &FeedTask{
		parser:   &SecondParser,
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
	next, err := t.parser.Parse(t.Cron())
	errors.NewAniErrorD(err).TryPanic()
	return next.Next(time.Now())
}

func (t *FeedTask) Run(args models.Object) {
	url := t.plugin.Get(VarUrl).(string)
	header := make(map[string]string)
	varHeader := t.plugin.Get(VarHeader)
	if varHeader != nil {
		for k, v := range varHeader.(map[string]any) {
			header[k] = v.(string)
		}
	}
	data, err := request.GetString(url, header)
	if err != nil {
		log.Warnf("[Plugin] %s插件(%s)执行错误: 请求 %s 失败", t.plugin.Type(), FuncParse, url)
		log.Debugf("", err)
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
			utils.MapToStruct(item.(map[string]any), items[i])
		}
		t.callback(items)
	}).Catch(func(err try.E) {
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", t.plugin.Type(), t.Name(), err)
		log.Debugf("", err)
	})

}

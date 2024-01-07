package schedule

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
)

type FeedTask struct {
	parser   *cron.Parser
	plugin   api.Plugin
	args     models.Object
	callback func([]*models.FeedItem) error
}

type FeedOptions struct {
	*models.Plugin
	Callback func([]*models.FeedItem) error
}

func NewFeedTask(opts *FeedOptions) (*FeedTask, error) {
	p, err := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    opts.Plugin,
		EntryFunc: FuncParse,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         FuncParse,
				ParamsSchema: []string{"data", "__retry_count__"},
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
	if err != nil {
		return nil, err
	}
	return &FeedTask{
		parser:   &SecondParser,
		plugin:   p,
		args:     opts.Args,
		callback: opts.Callback,
	}, nil
}

func (t *FeedTask) Name() string {
	name, err := t.plugin.Get(VarName)
	if err != nil {
		log.Warnf("%s", err)
	}
	if name == nil {
		name = "Feed"
	}
	return fmt.Sprintf("%v(%s-Plugin)", name, t.plugin.Type())
}

func (t *FeedTask) Cron() string {
	cronStr, err := t.plugin.Get(VarCron)
	if err != nil {
		log.Warnf("%s", err)
	}
	return cronStr.(string)
}

func (t *FeedTask) SetVars(vars models.Object) {
	for k, v := range vars {
		err := t.plugin.Set(k, v)
		if err != nil {
			log.Warnf("%s", err)
		}
	}
}

func (t *FeedTask) NextTime() time.Time {
	next, err := t.parser.Parse(t.Cron())
	if err != nil {
		log.DebugErr(err)
	}
	return next.Next(time.Now())
}

func (t *FeedTask) Run(args models.Object) (err error) {
	urlStr, err := t.plugin.Get(VarUrl)
	if err != nil {
		return err
	}
	url := urlStr.(string)
	header := make(map[string]string)
	varHeader, err := t.plugin.Get(VarHeader)
	if err != nil {
		return err
	}
	if varHeader != nil {
		for k, v := range varHeader.(map[string]any) {
			header[k] = v.(string)
		}
	}
	data, err := request.GetString(url, header)
	if err != nil {
		log.Warnf("[Plugin] %s插件(%s)执行错误: 请求 %s 失败", t.plugin.Type(), FuncParse, url)
		log.DebugErr(err)
	}
	args["data"] = data
	result, err := t.plugin.Run(FuncParse, args)
	if err != nil {
		return err
	}
	if result["error"] != nil {
		err = errors.WithStack(&exceptions.ErrScheduleRun{Name: t.Name(), Message: result["error"]})
		log.DebugErr(err)
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", t.plugin.Type(), t.Name(), result["error"])
		return err
	}
	itemsAny := result["items"].([]any)
	items := make([]*models.FeedItem, len(itemsAny))
	for i, item := range itemsAny {
		items[i] = &models.FeedItem{}
		err = utils.MapToStruct(item.(map[string]any), items[i])
		if err != nil {
			log.DebugErr(err)
			return errors.WithStack(&exceptions.ErrPlugin{Type: t.plugin.Type(), File: t.Name(), Message: "类型转换错误"})
		}
	}
	return t.callback(items)
}

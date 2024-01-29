package filter

import (
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
)

type Filter struct {
	plugin *models.Plugin
}

func NewFilterPlugin(pluginInfo *models.Plugin) *Filter {
	return &Filter{
		plugin: pluginInfo,
	}
}

func (p *Filter) FilterAll(items []*models.FeedItem) (resultItems []*models.FeedItem, err error) {
	if !p.plugin.Enable {
		err = errors.WithStack(&exceptions.ErrPluginDisabled{Type: p.plugin.Type, File: p.plugin.File})
		log.DebugErr(err)
		return nil, err
	}
	log.Debugf("[Plugin] 开始执行Filter插件(%s): %s", p.plugin.Type, p.plugin.File)
	// 入参
	pluginInstance, err := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    p.plugin,
		EntryFunc: constant.FuncFilterAll,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         constant.FuncFilterAll,
				ParamsSchema: []string{"items"},
				ResultSchema: []string{"error", "index"},
				DefaultArgs:  p.plugin.Args,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	result, err := pluginInstance.Run(constant.FuncFilterAll, map[string]any{
		"items": items,
	})
	if err != nil {
		return nil, err
	}
	if result["error"] != nil {
		err = errors.WithStack(&exceptions.ErrPlugin{Type: p.plugin.Type, File: p.plugin.File, Message: result["error"]})
		log.DebugErr(err)
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", p.plugin.Type, p.plugin.File, result["error"])
		return nil, err
	}

	if _, has := result["index"]; has {
		resultItems = filterIndex(items, result["index"].([]any))
	}
	return resultItems, nil
}

func filterIndex(items []*models.FeedItem, indexList []any) []*models.FeedItem {
	itemResult := make([]*models.FeedItem, len(indexList))
	for i, val := range indexList {
		index := int(val.(int64))
		if index < 0 || index >= len(items) {
			continue
		}
		itemResult[i] = items[index]
	}
	return itemResult
}

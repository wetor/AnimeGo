package plugin

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	FuncFilterAll = "filter_all"
)

type Filter struct {
	plugin *models.Plugin
}

func NewFilterPlugin(pluginInfo *models.Plugin) *Filter {
	return &Filter{
		plugin: pluginInfo,
	}
}

func (p *Filter) Filter(items []*models.FeedItem) (resultItems []*models.FeedItem) {

	if !p.plugin.Enable {
		return items
	}
	log.Debugf("[Plugin] 开始执行Filter插件(%s): %s", p.plugin.Type, p.plugin.File)
	// 入参
	pluginInstance := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    p.plugin,
		EntryFunc: FuncFilterAll,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         FuncFilterAll,
				ParamsSchema: []string{"items"},
				ResultSchema: []string{"error", "data,optional", "index,optional"},
				DefaultArgs:  p.plugin.Args,
			},
		},
	})
	result := pluginInstance.Run(FuncFilterAll, map[string]any{
		"items": items,
	})
	if result["error"] != nil {
		log.Debugf("", errors.NewAniErrorD(result["error"]))
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", p.plugin.Type, p.plugin.File, result["error"])
	}

	if _, has := result["data"]; has {
		resultItems = filterData(items, result["data"].([]any))
	} else if _, has := result["index"]; has {
		resultItems = filterIndex(items, result["index"].([]any))
	}
	return
}

func filterData(items []*models.FeedItem, data []any) []*models.FeedItem {
	itemResult := make([]*models.FeedItem, len(data))
	for i, val := range data {
		obj := val.(map[string]any)
		index := int(obj["index"].(int64))
		if index < 0 || index >= len(items) {
			continue
		}
		if _, has := obj["parsed"]; has {
			parsed := &models.TitleParsed{}
			utils.MapToStruct(obj["parsed"].(map[string]any), parsed)
			items[index].NameParsed = parsed
		}
		itemResult[i] = items[index]
	}
	return itemResult
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

package plugin

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

const (
	FuncFilterAll = "filter_all"
)

type Filter struct {
	plugin []models.Plugin
}

func NewFilterPlugin(pluginInfo []models.Plugin) *Filter {
	return &Filter{
		plugin: pluginInfo,
	}
}

func (p *Filter) Filter(list []*models.FeedItem) []*models.FeedItem {

	inList := make([]*models.FeedItem, len(list))
	for i, item := range list {
		inList[i] = item
	}
	for _, info := range p.plugin {
		if !info.Enable {
			continue
		}
		log.Debugf("[Plugin] 开始执行Filter插件(%s): %s", info.Type, info.File)
		// 入参
		pluginInstance := &python.Python{}
		pluginInstance.Load(&models.PluginLoadOptions{
			File: info.File,
			Functions: []*models.PluginFunctionOptions{
				{
					Name:         FuncFilterAll,
					ParamsSchema: []string{"items"},
					ResultSchema: []string{"error", "data,optional", "index,optional"},
				},
			},
		})
		result := pluginInstance.Run(FuncFilterAll, models.Object{
			"items": inList,
		})
		if result["error"] != nil {
			log.Debugf("", errors.NewAniErrorD(result["error"]))
			log.Warnf("[Plugin] %s插件(%s)执行错误: %v", info.Type, info.File, result["error"])
		}

		if _, has := result["data"]; has {
			inList = filterData(list, result["data"].([]any))
		} else if _, has := result["index"]; has {
			inList = filterIndex(list, result["index"].([]any))
		}
	}
	// 返回筛选结果
	return inList
}

func filterData(items []*models.FeedItem, data []any) []*models.FeedItem {
	itemResult := make([]*models.FeedItem, len(data))
	for i, val := range data {
		obj := val.(models.Object)
		index := int(obj["index"].(int64))
		if index < 0 || index >= len(items) {
			continue
		}
		if _, has := obj["parsed"]; has {
			parsed := &models.TitleParsed{}
			utils.MapToStruct(obj["parsed"].(models.Object), parsed)
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

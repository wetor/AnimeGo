package plugin

import (
	"path"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
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
		zap.S().Debugf("[Plugin] 开始执行Filter插件(%s): %s", info.Type, info.File)
		// 入参
		pluginInstance := plugin.GetPlugin(info.Type)
		pluginInstance.SetSchema([]string{"required:feedItems"}, []string{"required:error", "optional:data", "optional:index"})
		execute := pluginInstance.Execute(path.Join(constant.PluginPath, info.File), models.Object{
			"feedItems": inList,
		})
		result := execute.(models.Object)
		if result["error"] != nil {
			errors.NewAniErrorD(result["error"]).TryPanic()
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
			utils.Map2ModelByJson(obj["parsed"].(models.Object), parsed)
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
		if index < 0 || index >= len(indexList) {
			continue
		}
		itemResult[i] = items[index]
	}
	return itemResult
}

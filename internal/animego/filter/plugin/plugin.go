package plugin

import (
	"path"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
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
	// 过滤出的index列表
	filterIndex := make([]int64, 0, len(list))
	for i := range list {
		filterIndex = append(filterIndex, int64(i))
	}
	for _, info := range p.plugin {
		if !info.Enable {
			continue
		}
		zap.S().Debugf("[Plugin] 开始执行Filter插件(%s): %s", info.Type, info.File)
		// 入参
		inList := make([]*models.FeedItem, 0)
		for _, i := range filterIndex {
			inList = append(inList, list[i])
		}
		pluginInstance := plugin.GetPlugin(info.Type)
		pluginInstance.SetSchema([]string{"feedItems"}, []string{"index", "error"})
		execute := pluginInstance.Execute(path.Join(constant.PluginPath, info.File), models.Object{
			"feedItems": inList,
		})
		result := execute.(models.Object)
		if result["error"] != nil {
			errors.NewAniErrorD(result["error"]).TryPanic()
		}
		// 返回的index列表
		resultIndex := result["index"].([]any)

		filterIndex = make([]int64, 0, len(resultIndex))
		for _, index := range resultIndex {
			i := index.(int64)
			if i < 0 || i >= int64(len(list)) {
				continue
			}
			filterIndex = append(filterIndex, i)
		}
	}
	// 返回筛选结果
	result := make([]*models.FeedItem, 0, len(filterIndex))
	for _, index := range filterIndex {
		result = append(result, list[index])
	}
	return result
}

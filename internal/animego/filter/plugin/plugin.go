package plugin

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
)

type Filter struct {
	ScriptFile []string
	plugin     plugin.Plugin
}

func NewPluginFilter(plugin plugin.Plugin, files []string) *Filter {
	return &Filter{
		ScriptFile: files,
		plugin:     plugin,
	}
}

func (p *Filter) Filter(list []*models.FeedItem) []*models.FeedItem {
	// 过滤出的index列表
	filterIndex := make([]int64, 0, len(list))
	for i := range list {
		filterIndex = append(filterIndex, int64(i))
	}
	for _, jsFile := range p.ScriptFile {
		// 入参
		inList := make([]*models.FeedItem, 0)
		for _, i := range filterIndex {
			inList = append(inList, list[i])
		}
		p.plugin.SetSchema([]string{"feedItems"}, []string{"index", "error"})
		execute := p.plugin.Execute(jsFile, plugin.Object{
			"feedItems": inList,
		})
		// 返回的index列表
		resultIndex := execute.(plugin.Object)["index"].([]any)

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

package javascript

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/internal/store"
)

type JavaScript struct {
	ScriptFile []string
}

func (j *JavaScript) Filter(list []*models.FeedItem) []*models.FeedItem {
	if len(j.ScriptFile) == 0 {
		j.ScriptFile = store.Config.Filter.JavaScript
	}
	// 过滤出的index列表
	filterIndex := make([]int64, 0, len(list))
	for i := range list {
		filterIndex = append(filterIndex, int64(i))
	}
	for _, jsFile := range j.ScriptFile {
		// 入参
		inList := make([]*models.FeedItem, 0)
		for _, i := range filterIndex {
			inList = append(inList, list[i])
		}
		js := &javascript.JavaScript{}
		js.SetSchema([]string{"feedItems"}, []string{"index", "error"})
		execute := js.Execute(jsFile, javascript.Object{
			"feedItems": inList,
		})
		// 返回的index列表
		resultIndex := execute.(javascript.Object)["index"].([]any)

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

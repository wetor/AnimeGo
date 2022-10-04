package javascript

import (
	"AnimeGo/internal/animego/plugin/javascript"
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"go.uber.org/zap"
)

type JavaScript struct {
}

func (j *JavaScript) Filter(list []*models.FeedItem) []*models.FeedItem {
	js := &javascript.JavaScript{}
	js.SetSchema([]string{"feedItems"}, []string{"index", "error"})
	execute, err := js.Execute(store.Config.JavaScript, javascript.Object{
		"feedItems": list,
	})
	if err != nil {
		zap.S().Debug(err)
	}
	resultIndex := execute.(javascript.Object)["index"].([]interface{})
	result := make([]*models.FeedItem, 0, len(list))
	for _, index := range resultIndex {
		i := index.(int64)
		if i < 0 || i >= int64(len(list)) {
			continue
		}
		result = append(result, list[i])
	}
	return result
}

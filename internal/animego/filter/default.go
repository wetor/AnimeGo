package filter

import (
	"AnimeGo/internal/models"
)

type Default struct {
}

func (f *Default) Filter(list []*models.FeedItem) []*models.FeedItem {
	return list
}

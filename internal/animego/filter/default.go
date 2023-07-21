package filter

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type Default struct {
}

func (f *Default) FilterAll(list []*models.FeedItem) []*models.FeedItem {
	return list
}

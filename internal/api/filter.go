package api

import (
	"context"

	"github.com/wetor/AnimeGo/internal/models"
)

type FilterPlugin interface {
	FilterAll([]*models.FeedItem) ([]*models.FeedItem, error)
}

type FilterManager interface {
	Update(ctx context.Context, items []*models.FeedItem,
		parseOverride *models.AnimeParseOverride,
		skipFilter bool, skipDelay bool) error
}

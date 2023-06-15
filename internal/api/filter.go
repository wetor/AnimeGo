package api

import (
	"context"

	"github.com/wetor/AnimeGo/internal/models"
)

type FilterPlugin interface {
	FilterAll([]*models.FeedItem) []*models.FeedItem
}

type FilterManager interface {
	Update(ctx context.Context, items []*models.FeedItem, parseOverride *models.AnimeParseOverride, skipFilter bool)
}

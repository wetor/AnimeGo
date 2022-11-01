package process

import (
	"context"
	"github.com/wetor/AnimeGo/internal/models"
)

type Process interface {
	UpdateFeed(items []*models.FeedItem)
	Run(ctx context.Context)
}

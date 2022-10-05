package process

import (
	"AnimeGo/internal/models"
	"context"
)

type Process interface {
	UpdateFeed(items []*models.FeedItem)
	Run(ctx context.Context)
}

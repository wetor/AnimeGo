package api

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type Feed interface {
	Parse(opts ...any) []*models.FeedItem
}

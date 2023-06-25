package api

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type Feed interface {
	Parse() ([]*models.FeedItem, error)
}

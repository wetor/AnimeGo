package api

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type Feed interface {
	ParseFile(string) ([]*models.FeedItem, error)
	ParseUrl(string) ([]*models.FeedItem, error)
	Parse([]byte) ([]*models.FeedItem, error)
}

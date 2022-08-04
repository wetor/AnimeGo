package feed

import (
	"GoBangumi/internal/models"
)

type Feed interface {
	Parse(opt *models.FeedParseOptions) []*models.FeedItem
}

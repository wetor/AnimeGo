package feed

import (
	"GoBangumi/models"
)

type Feed interface {
	Parse(opt *models.FeedParseOptions) []*models.FeedItem
}

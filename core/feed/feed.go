package feed

import (
	"GoBangumi/model"
)

type Feed interface {
	Parse(opt *model.FeedParseOptions) []*model.FeedItem
}

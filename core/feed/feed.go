package feed

import (
	"GoBangumi/core/bangumi"
	"GoBangumi/model"
	"github.com/mmcdole/gofeed"
)

type Feed interface {
	Parse(opt *model.FeedParseOptions)
	ParseBangumiAll(bangumi bangumi.Bangumi) []*model.Bangumi
	ParseBangumi(item *gofeed.Item, bangumi bangumi.Bangumi) *model.Bangumi
}

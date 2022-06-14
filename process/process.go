package process

import (
	"GoBangumi/core/bangumi"
	"GoBangumi/model"
)

type Process interface {
	Run()
	ParseBangumiAll(items []*model.FeedItem, bangumi bangumi.Bangumi) []*model.Bangumi
	ParseBangumi(item *model.FeedItem, bangumi bangumi.Bangumi) *model.Bangumi
}

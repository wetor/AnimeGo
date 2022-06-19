package process

import (
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
)

type Process interface {
	Run()
	ParseBangumiAll(items []*models.FeedItem, bangumi bangumi.Bangumi) []*models.Bangumi
	ParseBangumi(item *models.FeedItem, bangumi bangumi.Bangumi) *models.Bangumi
}

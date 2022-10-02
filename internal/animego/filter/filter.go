// Package filter
// @Description: 过滤器包，用来过滤符合条件的下载条目
package filter

import "AnimeGo/internal/models"

type Filter interface {
	Filter([]*models.FeedItem) []*models.FeedItem
}

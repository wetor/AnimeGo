// Package filter
// @Description: 过滤器包，用来过滤符合条件的下载条目
package filter

import "GoBangumi/internal/models"

type Filter interface {
	Filter([]*models.FeedItem) []*models.FeedItem
}

func ResolutionFilter(src, filter string) bool {
	// 720, 1080, 4k
	return true
}

func SubtitleFilter(src, filter string) bool {
	// 简体, 繁体
	return true
}

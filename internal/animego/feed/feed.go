// Package feed
// @Description: 订阅包，主要用来解析订阅信息
package feed

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type Feed interface {
	Parse(opts ...any) []*models.FeedItem
}

var (
	TempPath string
)

type Options struct {
	TempPath string
}

func Init(opts *Options) {
	TempPath = opts.TempPath
}

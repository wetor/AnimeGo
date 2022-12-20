// Package downloader
// @Description: 下载器包，用来调用外部下载器
package downloader

import (
	"context"
	"github.com/wetor/AnimeGo/internal/models"
)

type Client interface {
	Connected() bool
	Start(ctx context.Context)
	List(opt *models.ClientListOptions) []*models.TorrentItem
	Add(opt *models.ClientAddOptions)
	Delete(opt *models.ClientDeleteOptions)
	GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem
}

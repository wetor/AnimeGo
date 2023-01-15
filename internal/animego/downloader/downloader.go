// Package downloader
// @Description: 下载器包，用来调用外部下载器
package downloader

import (
	"context"
	"sync"

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

var (
	ConnectTimeoutSecond int
	CheckTimeSecond      int
	RetryConnectNum      int
	WG                   *sync.WaitGroup
)

type Options struct {
	ConnectTimeoutSecond int
	CheckTimeSecond      int
	RetryConnectNum      int
	WG                   *sync.WaitGroup
}

func Init(opts *Options) {
	ConnectTimeoutSecond = opts.ConnectTimeoutSecond
	CheckTimeSecond = opts.CheckTimeSecond
	RetryConnectNum = opts.RetryConnectNum
	WG = opts.WG
}

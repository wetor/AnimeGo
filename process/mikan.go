package process

import (
	"GoBangumi/internal/core/anisource/mikan"
	feedManager "GoBangumi/internal/core/feed/manager"
	mikanRss "GoBangumi/internal/core/feed/mikan"
	downloaderManager "GoBangumi/internal/core/manager"
	"GoBangumi/internal/downloader/qbittorent"
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"context"
)

const (
	UpdateWaitMinMinute = 2 // 订阅最短间隔分钟
)

type Mikan struct {
	downloaderMgr *downloaderManager.Manager
	feedMgr       *feedManager.Manager
	exitChan      chan bool // 结束标记
}

func NewMikan() *Mikan {
	return &Mikan{
		exitChan: make(chan bool),
	}
}
func (p *Mikan) Exit() {
	p.exitChan <- true
}
func (p *Mikan) Run(ctx context.Context) {

	qbtConf := store.Config.ClientQBt()
	qbt := qbittorent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password)

	downloadChan := make(chan *models.AnimeEntity, 10)

	p.downloaderMgr = downloaderManager.NewManager(qbt, downloadChan)
	p.feedMgr = feedManager.NewManager(mikanRss.NewRss(), mikan.NewMikan())
	p.feedMgr.SetDownloadChan(downloadChan)

	p.downloaderMgr.Start(ctx)
	p.feedMgr.Start(ctx)
}

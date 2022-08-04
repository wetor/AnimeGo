package process

import (
	"GoBangumi/internal/core/anisource/mikan"
	feedManager "GoBangumi/internal/core/feed/manager"
	mikanRss "GoBangumi/internal/core/feed/mikan"
	downloaderManager "GoBangumi/internal/downloader/manager"
	"GoBangumi/internal/downloader/qbittorent"
	"GoBangumi/internal/models"
	"GoBangumi/store"
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
	return &Mikan{}
}
func (p *Mikan) Exit() {
	p.exitChan <- true
}
func (p *Mikan) Run(exit chan bool) {

	qbtConf := store.Config.ClientQBt()
	qbt := qbittorent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password)

	downloadChan := make(chan *models.AnimeEntity, 10)

	p.downloaderMgr = downloaderManager.NewManager(qbt)
	p.downloaderMgr.SetDownloadChan(downloadChan)
	p.feedMgr = feedManager.NewManager(mikanRss.NewRss(), mikan.NewMikan())
	p.feedMgr.SetDownloadChan(downloadChan)

	downloaderExit := make(chan bool)
	feedExit := make(chan bool)
	p.downloaderMgr.Start(downloaderExit)
	p.feedMgr.Start(feedExit)

	exitFlag := 0
	go func() {
		select {
		case <-feedExit:
			exitFlag++
		case <-downloaderExit:
			exitFlag++
		default:
			if exitFlag >= 2 {
				exit <- true
				return
			}
		}
	}()
}

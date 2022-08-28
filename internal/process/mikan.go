package process

import (
	"GoBangumi/internal/animego/anisource"
	"GoBangumi/internal/animego/anisource/mikan"
	"GoBangumi/internal/animego/downloader/qbittorent"
	mikanRss "GoBangumi/internal/animego/feed/mikan"
	"GoBangumi/internal/animego/filter"
	"GoBangumi/internal/animego/manager/downloader"
	filterManager "GoBangumi/internal/animego/manager/filter"
	"GoBangumi/internal/models"
	"GoBangumi/internal/store"
	"context"
)

const (
	UpdateWaitMinMinute = 2 // 订阅最短间隔分钟
)

type Mikan struct {
	downloaderMgr *downloader.Manager
	filterMgr     *filterManager.Manager
}

func NewMikan() *Mikan {
	return &Mikan{}
}

func (p *Mikan) Run(ctx context.Context) {

	qbtConf := store.Config.ClientQBt()
	qbt := qbittorent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password, ctx)

	downloadChan := make(chan *models.AnimeEntity, 10)

	anisource.Init(store.Cache, store.Config.Proxy())

	p.downloaderMgr = downloader.NewManager(qbt, store.Cache, downloadChan)

	p.filterMgr = filterManager.NewManager(&filter.Default{},
		mikanRss.NewRss(store.Config.RssMikan().Url, store.Config.RssMikan().Name),
		mikan.MikanAdapter{ThemoviedbKey: store.Config.KeyTmdb()},
		downloadChan)

	p.downloaderMgr.Start(ctx)
	p.filterMgr.Start(ctx)
}

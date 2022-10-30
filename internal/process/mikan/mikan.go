package mikan

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/animego/anisource/mikan"
	"AnimeGo/internal/animego/downloader/qbittorent"
	mikanRss "AnimeGo/internal/animego/feed/mikan"
	"AnimeGo/internal/animego/filter/javascript"
	"AnimeGo/internal/animego/manager/downloader"
	filterManager "AnimeGo/internal/animego/manager/filter"
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	pkgAnisource "AnimeGo/pkg/anisource"
	"context"
)

type Mikan struct {
	downloaderMgr *downloader.Manager
	filterMgr     *filterManager.Manager
	ctx           context.Context
}

func NewMikan() *Mikan {
	return &Mikan{}
}

func (p *Mikan) UpdateFeed(items []*models.FeedItem) {
	p.filterMgr.Update(p.ctx, items)
}

func (p *Mikan) Run(ctx context.Context) {

	qbtConf := store.Config.Setting.Client.QBittorrent
	qbt := qbittorent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password)
	qbt.Start(ctx)

	downloadChan := make(chan *models.AnimeEntity, 10)

	anisource.Init(&pkgAnisource.Options{
		Cache: store.Cache,
	})

	p.downloaderMgr = downloader.NewManager(qbt, store.Cache, downloadChan)

	p.filterMgr = filterManager.NewManager(&javascript.JavaScript{},
		mikanRss.NewRss(store.Config.Setting.Feed.Mikan.Url, store.Config.Setting.Feed.Mikan.Name),
		mikan.MikanAdapter{ThemoviedbKey: store.Config.Setting.Key.Themoviedb},
		downloadChan)
	p.ctx = ctx
	p.downloaderMgr.Start(ctx)
	p.filterMgr.Start(ctx)
}

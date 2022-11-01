package mikan

import (
	"context"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorent"
	mikanRss "github.com/wetor/AnimeGo/internal/animego/feed/mikan"
	"github.com/wetor/AnimeGo/internal/animego/filter/javascript"
	"github.com/wetor/AnimeGo/internal/animego/manager/downloader"
	filterManager "github.com/wetor/AnimeGo/internal/animego/manager/filter"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	pkgAnisource "github.com/wetor/AnimeGo/pkg/anisource"
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

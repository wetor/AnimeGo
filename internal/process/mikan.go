package process

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
}

func NewMikan() *Mikan {
	return &Mikan{}
}

func (p *Mikan) Run(ctx context.Context) {

	qbtConf := store.Config.ClientQBt()
	qbt := qbittorent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password, ctx)

	downloadChan := make(chan *models.AnimeEntity, 10)

	anisource.Init(&pkgAnisource.Options{
		Cache:   store.Cache,
		Proxy:   store.Config.Proxy(),
		Timeout: store.Config.HttpTimeoutSecond,
		Retry:   store.Config.HttpRetryNum,
	})

	p.downloaderMgr = downloader.NewManager(qbt, store.Cache, downloadChan)

	p.filterMgr = filterManager.NewManager(&javascript.JavaScript{},
		mikanRss.NewRss(store.Config.RssMikan().Url, store.Config.RssMikan().Name),
		mikan.MikanAdapter{ThemoviedbKey: store.Config.KeyTmdb()},
		downloadChan)

	p.downloaderMgr.Start(ctx)
	p.filterMgr.Start(ctx)
}

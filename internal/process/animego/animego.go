package animego

import (
	"context"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	anidataBangumi "github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	anidataThemoviedb "github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	feedRss "github.com/wetor/AnimeGo/internal/animego/feed/rss"
	"github.com/wetor/AnimeGo/internal/animego/filter/plugin"
	"github.com/wetor/AnimeGo/internal/animego/manager/downloader"
	filterManager "github.com/wetor/AnimeGo/internal/animego/manager/filter"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/internal/schedule"
	"github.com/wetor/AnimeGo/internal/store"
)

type AnimeGo struct {
	downloaderMgr *downloader.Manager
	filterMgr     *filterManager.Manager
	schedule      *schedule.Schedule
	ctx           context.Context
}

func NewAnimeGo() *AnimeGo {
	return &AnimeGo{}
}

func (p *AnimeGo) Update(data any) {
	items := data.([]*models.FeedItem)
	p.filterMgr.Update(p.ctx, items)
}

func (p *AnimeGo) Run(ctx context.Context) {

	p.schedule = schedule.NewSchedule()
	p.schedule.Start(ctx)

	qbtConf := store.Config.Setting.Client.QBittorrent
	qbt := qbittorrent.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password)
	qbt.Start(ctx)

	downloadChan := make(chan *models.AnimeEntity, 10)

	anisource.Init(&anidata.Options{
		Cache: store.Cache,
		CacheTime: map[string]int64{
			anidataBangumi.Bucket:    int64(store.Config.Advanced.Cache.MikanCacheHour * 60 * 60),
			anidataBangumi.Bucket:    int64(store.Config.Advanced.Cache.BangumiCacheHour * 60 * 60),
			anidataThemoviedb.Bucket: int64(store.Config.Advanced.Cache.ThemoviedbCacheHour * 60 * 60),
		},
	})

	p.downloaderMgr = downloader.NewManager(qbt, store.Cache, downloadChan)

	p.filterMgr = filterManager.NewManager(
		plugin.NewPluginFilter(&javascript.JavaScript{}, store.Config.Filter.JavaScript),
		feedRss.NewRss(store.Config.Setting.Feed.Mikan.Url, store.Config.Setting.Feed.Mikan.Name),
		mikan.MikanAdapter{ThemoviedbKey: store.Config.Setting.Key.Themoviedb},
		downloadChan)
	p.ctx = ctx
	p.downloaderMgr.Start(ctx)
	p.filterMgr.Start(ctx)
}

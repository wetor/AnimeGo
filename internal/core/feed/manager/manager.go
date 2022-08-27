package manager

import (
	"GoBangumi/internal/core/anisource"
	"GoBangumi/internal/core/feed"
	mikanRss "GoBangumi/internal/core/feed/mikan"
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"GoBangumi/utils"
	"context"
	"go.uber.org/zap"
	"sync"
)

const (
	UpdateWaitMinMinute = 2 // 订阅最短间隔分钟
)

type Manager struct {
	feed      feed.Feed
	anisource anisource.AniSource
	animeList []*models.AnimeEntity

	downloadChanEnable bool
	downloadChan       chan *models.AnimeEntity
}

// NewManager
//  @Description:
//  @param feed feed.Feed
//  @param anisource1 anisource1.AniSource
//  @return *Manager
//
func NewManager(feed feed.Feed, anisource anisource.AniSource) *Manager {
	m := &Manager{
		feed:      feed,
		anisource: anisource,
	}
	return m
}

func (m *Manager) SetDownloadChan(donwloadChan chan *models.AnimeEntity) {
	if cap(donwloadChan) > 0 {
		m.downloadChanEnable = true
		m.downloadChan = donwloadChan
	}
}

func (m *Manager) GetAnimeList() []*models.AnimeEntity {
	list := make([]*models.AnimeEntity, len(m.animeList))
	copy(list, m.animeList)
	return list
}

func (m *Manager) UpdateFeed(ctx context.Context) {
	rssConf := store.Config.RssMikan()
	f := mikanRss.NewRss()
	items := f.Parse(&models.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name, // 文件名
		RefreshCache: true,
	})

	conf := store.Config.Advanced.MainConf
	animeList := make([]*models.AnimeEntity, len(items))
	working := make(chan int, conf.MultiGoroutine.GoroutineMax) // 限制同时执行个数
	wg := sync.WaitGroup{}
	exit := false
	for i, item := range items {
		if exit {
			return
		}
		working <- i //计数器+1 可能会发生阻塞
		wg.Add(1)
		go func(_i int, _item *models.FeedItem) {
			defer func() {
				if err := recover(); err != nil {
					zap.S().Error(err)
				}
			}()
			select {
			case <-ctx.Done():
				exit = true
			default:
				anime := m.anisource.Parse(&models.AnimeParseOptions{
					Url:  _item.Url,
					Name: _item.Name,
					Date: _item.Date,
				})
				if anime != nil {
					if anime.TorrentInfo == nil {
						anime.TorrentInfo = &models.TorrentInfo{}
					}
					anime.Url = _item.Torrent
					anime.Hash = _item.Hash

					animeList[_i] = anime
					if m.downloadChanEnable {
						zap.S().Debugf("发送下载项:「%s」", anime.FullName())
						// 向管道中发送需要下载的信息
						m.downloadChan <- anime
					}
				}
				utils.Sleep(conf.FeedDelay, ctx)
			}
			<-working
			wg.Done()
		}(i, item)

		if !exit && !conf.MultiGoroutine.Enable {
			wg.Wait()
		}
	}
	// 等待处理完成
	wg.Wait()
}

func (m *Manager) Start(ctx context.Context) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				zap.S().Error(err)
			}
		}()
		defer store.WG.Done()
		for {
			select {
			case <-ctx.Done():
				zap.S().Debug("正常退出")
				return
			default:
				m.UpdateFeed(ctx)
				delay := store.Config.Advanced.MainConf.FeedUpdateDelayMinute
				if delay < UpdateWaitMinMinute {
					delay = UpdateWaitMinMinute
				}
				utils.Sleep(delay*60, ctx)
			}
		}
	}()
}

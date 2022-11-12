// Package filter
// @Description: 筛选输入feed条目，并通过anisource获取符合条目的详细信息，信息完整则传递给下载器进行下载
package filter

import (
	"context"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"go.uber.org/zap"
	"sync"
)

const (
	UpdateWaitMinMinute    = 2  // 订阅最短间隔分钟
	DownloadChanDefaultCap = 10 // 下载通道默认容量
)

type Manager struct {
	filter       filter.Filter
	feed         feed.Feed
	anisource    anisource.AniSource
	downloadChan chan *models.AnimeEntity

	animeList []*models.AnimeEntity
}

// NewManager
//  @Description:
//  @param filter *filter.Filter
//  @param feed feed.Feed
//  @param anisource anisource.AniSource
//  @return *Manager
//
func NewManager(filter filter.Filter, feed feed.Feed, anisource anisource.AniSource, downloadChan chan *models.AnimeEntity) *Manager {
	m := &Manager{
		filter:    filter,
		feed:      feed,
		anisource: anisource,
	}
	if downloadChan == nil || cap(downloadChan) <= 1 {
		downloadChan = make(chan *models.AnimeEntity, DownloadChanDefaultCap)
	}
	m.downloadChan = downloadChan
	return m
}

func (m *Manager) Update(ctx context.Context, items []*models.FeedItem) {
	// 筛选
	if items == nil {
		items, _ = m.feed.Parse()
	}
	items = m.filter.Filter(items)

	animeList := make([]*models.AnimeEntity, len(items))
	working := make(chan int, store.Config.Advanced.Feed.MultiGoroutine.GoroutineMax) // 限制同时执行个数
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
					if anime.DownloadInfo == nil {
						anime.DownloadInfo = &models.DownloadInfo{}
					}
					anime.Url = _item.Download
					anime.Hash = _item.Hash()

					animeList[_i] = anime

					zap.S().Debugf("发送下载项:「%s」", anime.FullName())
					// 向管道中发送需要下载的信息
					m.downloadChan <- anime

				}
				utils.Sleep(store.Config.Advanced.Feed.Delay, ctx)
			}
			<-working
			wg.Done()
		}(i, item)

		if !exit && !store.Config.Advanced.Feed.MultiGoroutine.Enable {
			wg.Wait()
		}
	}
	// 等待处理完成
	wg.Wait()
}

func (m *Manager) Start(ctx context.Context) {
	store.WG.Add(1)
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
				zap.S().Debug("正常退出 manager filter")
				return
			default:
				m.Update(ctx, nil)
				delay := store.Config.Advanced.Feed.UpdateDelayMinute
				if delay < UpdateWaitMinMinute {
					delay = UpdateWaitMinMinute
				}
				utils.Sleep(delay*60, ctx)
			}
		}
	}()
}

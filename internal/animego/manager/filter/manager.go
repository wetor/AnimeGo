// Package filter
// @Description: 筛选输入feed条目，并通过anisource获取符合条目的详细信息，信息完整则传递给下载器进行下载
package filter

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
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
//
//	@Description:
//	@param filter *filter.Filter
//	@param feed feed.Feed
//	@param anisource anisource.AniSource
//	@return *Manager
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
		items = m.feed.Parse()
	}
	if len(items) == 0 {
		return
	}
	items = m.filter.Filter(items)

	animeList := make([]*models.AnimeEntity, len(items))
	working := make(chan int, manager.FilterConf.MultiGoroutineMax) // 限制同时执行个数
	wg := sync.WaitGroup{}
	exit := false
	for i, item := range items {
		if exit {
			return
		}
		working <- i //计数器+1 可能会发生阻塞
		wg.Add(1)
		go func(_i int, _item *models.FeedItem) {
			defer errors.HandleError(func(err error) {
				zap.S().Error(err)
			})
			select {
			case <-ctx.Done():
				exit = true
			default:
				anime := m.anisource.Parse(&models.AnimeParseOptions{
					Url:    _item.Url,
					Name:   _item.Name,
					Parsed: _item.NameParsed,
				})
				if anime != nil {
					anime.DownloadInfo = &models.DownloadInfo{
						Url:  _item.Download,
						Hash: _item.Hash(),
					}

					animeList[_i] = anime

					zap.S().Debugf("发送下载项:「%s」", anime.FullName())
					// 向管道中发送需要下载的信息
					m.downloadChan <- anime

				}
				utils.Sleep(manager.FilterConf.DelaySecond, ctx)
			}
			<-working
			wg.Done()
		}(i, item)

		if !exit && !manager.FilterConf.MultiGoroutineEnabled {
			wg.Wait()
		}
	}
	// 等待处理完成
	wg.Wait()
}

func (m *Manager) Start(ctx context.Context) {
	manager.WG.Add(1)
	go func() {
		defer manager.WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					zap.S().Error(err)
				})
				select {
				case <-ctx.Done():
					zap.S().Debug("正常退出 manager filter")
					exit = true
					return
				default:
					m.Update(ctx, nil)
					delay := manager.FilterConf.UpdateDelayMinute
					if delay < UpdateWaitMinMinute {
						delay = UpdateWaitMinMinute
					}
					utils.Sleep(delay*60, ctx)
				}
			}()
			if exit {
				return
			}
		}
	}()
}

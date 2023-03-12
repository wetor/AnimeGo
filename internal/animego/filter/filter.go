// Package filter
// @Description: 筛选输入feed条目，并通过anisource获取符合条目的详细信息，信息完整则传递给下载器进行下载
package filter

import (
	"context"
	"sync"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	DownloadChanDefaultCap = 10 // 下载通道默认容量
)

type Manager struct {
	filter       api.Filter
	anisource    api.AniSource
	downloadChan chan *models.AnimeEntity
	animeList    []*models.AnimeEntity
}

// NewFilter
//
//	@Description:
//	@param filter *filter.Filter
//	@param feed api.Feed
//	@param anisource api.AniSource
//	@return *Manager
func NewFilter(filter api.Filter, anisource api.AniSource, downloadChan chan *models.AnimeEntity) *Manager {
	m := &Manager{
		filter:    filter,
		anisource: anisource,
	}
	if downloadChan == nil || cap(downloadChan) <= 1 {
		downloadChan = make(chan *models.AnimeEntity, DownloadChanDefaultCap)
	}
	m.downloadChan = downloadChan
	return m
}

func (m *Manager) Update(ctx context.Context, items []*models.FeedItem) {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	// 筛选
	if len(items) == 0 {
		return
	}
	items = m.filter.Filter(items)

	animeList := make([]*models.AnimeEntity, len(items))
	working := make(chan int, MultiGoroutineMax) // 限制同时执行个数
	wg := sync.WaitGroup{}
	exit := false
	for i, item := range items {
		if exit {
			return
		}
		working <- i //计数器+1 可能会发生阻塞
		wg.Add(1)
		go func(_i int, _item *models.FeedItem) {
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

					log.Debugf("发送下载项:「%s」", anime.FullName())
					// 向管道中发送需要下载的信息
					//m.downloadChan <- anime

				}
				utils.Sleep(DelaySecond, ctx)
			}
			<-working
			wg.Done()
		}(i, item)

		if !exit && !MultiGoroutineEnabled {
			wg.Wait()
		}
	}
	// 等待处理完成
	wg.Wait()
}

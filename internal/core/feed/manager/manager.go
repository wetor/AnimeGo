package manager

import (
	"GoBangumi/internal/core/anisource"
	"GoBangumi/internal/core/feed"
	mikanRss "GoBangumi/internal/core/feed/mikan"
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"GoBangumi/utils"
	"go.uber.org/zap"
	"sync"
	"time"
)

const (
	UpdateWaitMinMinute = 2 // 订阅最短间隔分钟
)

type Manager struct {
	feed      feed.Feed
	anisource anisource.AniSource
	animeList []*models.AnimeEntity

	exitChan chan bool // 结束标记

	downloadChanEnable bool
	downloadChan       chan *models.AnimeEntity
}

// NewManager
//  @Description:
//  @param feed feed.Feed
//  @param anisource anisource.AniSource
//  @return *Manager
//
func NewManager(feed feed.Feed, anisource anisource.AniSource) *Manager {
	m := &Manager{
		feed:      feed,
		anisource: anisource,
		exitChan:  make(chan bool),
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

func (m *Manager) UpdateFeed() {
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
	for i, item := range items {
		working <- i //计数器+1 可能会发生阻塞
		wg.Add(1)
		go func(_i int, _item *models.FeedItem) {
			anime := m.anisource.Parse(&models.AnimeParseOptions{
				Url:  _item.Url,
				Name: _item.Name,
				Date: _item.Date,
			})
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
			time.Sleep(time.Duration(conf.FeedDelay) * time.Second)

			//工作完成后计数器减1
			<-working
			wg.Done()

		}(i, item)
		if !conf.MultiGoroutine.Enable {
			wg.Wait()
		}
	}
	// 等待处理完成
	wg.Wait()
}

func (m *Manager) Exit() {
	m.exitChan <- true
}
func (m *Manager) Start(exit chan bool) {
	go func() {
		select {
		case <-m.exitChan:
			exit <- true
			return
		default:
			m.UpdateFeed()
			delay := store.Config.Advanced.MainConf.FeedUpdateDelayMinute
			if delay < UpdateWaitMinMinute {
				delay = UpdateWaitMinMinute
			}
			if utils.Sleep(delay*60, m.exitChan) {
				exit <- true
				return
			}
		}
	}()
}

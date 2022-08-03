package process

import (
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
	"GoBangumi/modules/client"
	"GoBangumi/modules/feed"
	"GoBangumi/process/manager"
	"GoBangumi/store"
	"sync"
	"time"
)

const (
	UpdateWaitMinMinute = 2 // 订阅最短间隔分钟
)

type Mikan struct {
	mgr      *manager.Manager
	exitChan chan bool // 结束标记
}

func NewMikan() *Mikan {
	return &Mikan{}
}
func (p *Mikan) Exit() {
	p.exitChan <- true
}
func (p *Mikan) Run(exit chan bool) {

	qbtConf := store.Config.ClientQBt()
	qbt := client.NewQBittorrent(qbtConf.Url, qbtConf.Username, qbtConf.Password)
	p.mgr = manager.NewManager(qbt)
	managerExit := make(chan bool)
	p.mgr.Start(managerExit)

	go func() {
		select {
		case <-p.exitChan:
			p.mgr.Exit()
			// 等待manager退出
			<-managerExit
			// exit02
			exit <- true
			return
		default:
			p.UpdateFeed()
			delay := store.Config.FeedUpdateDelayMinute
			if delay < UpdateWaitMinMinute {
				delay = UpdateWaitMinMinute
			}
			time.Sleep(time.Duration(delay) * time.Minute)
		}
	}()
}

func (p *Mikan) UpdateFeed() {
	rssConf := store.Config.RssMikan()
	f := feed.NewRss()
	items := f.Parse(&models.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name, // 文件名
		RefreshCache: true,
	})
	bgms := p.ParseBangumiAll(items, &bangumi.Mikan{})

	for _, bgm := range bgms {
		p.mgr.Download(bgm)
	}
}

func (p *Mikan) ParseBangumiAll(items []*models.FeedItem, bangumi bangumi.Bangumi) []*models.Bangumi {
	bgms := make([]*models.Bangumi, len(items))
	conf := store.Config.Advanced.MainConf
	working := make(chan int, conf.MultiGoroutine.GoroutineMax) // 限制同时执行个数
	wg := sync.WaitGroup{}
	for i, item := range items {
		working <- i //计数器+1 可能会发生阻塞
		wg.Add(1)
		go func(i_ int, item_ *models.FeedItem) {
			defer wg.Done()
			bgms[i_] = p.ParseBangumi(item_, bangumi)
			if bgms[i_].TorrentInfo == nil {
				bgms[i_].TorrentInfo = &models.TorrentInfo{}
			}
			bgms[i_].Url = item_.Torrent
			bgms[i_].Hash = item_.Hash
			time.Sleep(time.Duration(conf.FeedDelay) * time.Second)
			//工作完成后计数器减1
			<-working
		}(i, item)
		if !conf.MultiGoroutine.Enable {
			wg.Wait() // 同步
		}
	}
	wg.Wait()
	return bgms
}

func (p *Mikan) ParseBangumi(item *models.FeedItem, bangumi bangumi.Bangumi) *models.Bangumi {
	bgmInfo := bangumi.Parse(&models.BangumiParseOptions{
		Url:  item.Url,
		Name: item.Name,
		Date: item.Date,
	})
	return bgmInfo
}

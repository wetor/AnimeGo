package process

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
	"GoBangumi/modules/feed"
	"fmt"
	"sync"
	"time"
)

type Mikan struct {
}

func NewMikan() Process {
	return &Mikan{}
}

func (p *Mikan) Run() {
	rssConf := config.RssMikan()

	f := feed.NewRss()
	items := f.Parse(&models.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name, // 文件名
		RefreshCache: true,
	})
	bgms := p.ParseBangumiAll(items, &bangumi.Mikan{})
	fmt.Println(bgms)
	for _, b := range bgms {
		fmt.Println(b)
	}
}

func (p *Mikan) ParseBangumiAll(items []*models.FeedItem, bangumi bangumi.Bangumi) []*models.Bangumi {
	bgms := make([]*models.Bangumi, len(items))
	conf := config.Advanced().GoBangumi()
	working := make(chan int, conf.MultiGoroutine.GoroutineMax) // 限制同时执行个数
	wg := sync.WaitGroup{}
	for i, item := range items {
		working <- i //计数器+1 可能会发生阻塞
		wg.Add(1)
		go func(i_ int, item_ *models.FeedItem) {
			defer wg.Done()
			bgms[i_] = p.ParseBangumi(item_, bangumi)
			if bgms[i_].BangumiExtra == nil {
				bgms[i_].BangumiExtra = &models.BangumiExtra{}
			}
			bgms[i_].TorrentUrl = item_.Torrent
			bgms[i_].TorrentHash = item_.Hash
			time.Sleep(time.Duration(conf.RssDelay) * time.Second)
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
	// TODO: 读取缓存，若存在且不过期则直接返回
	bgmInfo := bangumi.Parse(&models.BangumiParseOptions{
		Url:  item.Url,
		Name: item.Name,
		Date: item.Date,
	})
	// TODO: 写入缓存，需要线程安全
	return bgmInfo
}

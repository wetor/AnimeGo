package process

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
	"GoBangumi/modules/feed"
	"fmt"
	"sync"
)

type Mikan struct {
}

func NewMikan() Process {
	return &Mikan{}
}

func (p *Mikan) Run() {
	rssConf := config.Mikan()

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
	wg := sync.WaitGroup{}
	for i, item := range items {
		wg.Add(1)
		go func(i_ int, item_ *models.FeedItem) {
			defer wg.Done()
			bgms[i_] = p.ParseBangumi(item_, bangumi)
		}(i, item)
		wg.Wait() // 同步
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

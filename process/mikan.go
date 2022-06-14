package process

import (
	"GoBangumi/config"
	"GoBangumi/core/bangumi"
	"GoBangumi/core/feed"
	"GoBangumi/model"
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
	items := f.Parse(&model.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name,
		RefreshCache: true,
	})
	bgms := p.ParseBangumiAll(items, &bangumi.Mikan{})
	fmt.Println(bgms)
	for _, b := range bgms {
		fmt.Println(b)
	}
}

func (p *Mikan) ParseBangumiAll(items []*model.FeedItem, bangumi bangumi.Bangumi) []*model.Bangumi {
	bgms := make([]*model.Bangumi, len(items))
	wg := sync.WaitGroup{}
	for i, item := range items {
		wg.Add(1)
		go func(i_ int, item_ *model.FeedItem) {
			defer wg.Done()
			bgms[i_] = p.ParseBangumi(item_, bangumi)
		}(i, item)
		wg.Wait() // 同步
	}
	wg.Wait()
	return bgms
}
func (p *Mikan) ParseBangumi(item *model.FeedItem, bangumi bangumi.Bangumi) *model.Bangumi {
	// TODO: 读取缓存，若存在且不过期则直接返回
	bgmInfo := bangumi.Parse(&model.BangumiParseOptions{
		Url:  item.Url,
		Name: item.Name,
	})
	// TODO: 写入缓存，需要线程安全
	return bgmInfo
}

package feed

import (
	"GoBangumi/config"
	"GoBangumi/core/bangumi"
	"GoBangumi/model"
	"GoBangumi/utils"
	"github.com/golang/glog"
	"github.com/mmcdole/gofeed"
	"os"
	"path"
	"sync"
)

type Rss struct {
	url  string
	name string
	file string
	feed *gofeed.Feed
}

func NewRss() Feed {
	return &Rss{}
}

// Parse
//  Description 第一步，解析rss
//  Receiver f *Rss
//  Param opt *model.FeedParseOptions 若RefreshCache为false，则仅重新解析本地缓存rss
//
func (f *Rss) Parse(opt *model.FeedParseOptions) {
	// --------- 是否刷新信息 ---------
	if len(opt.Url) != 0 {
		f.url = opt.Url
	}
	if len(opt.Name) != 0 {
		f.name = opt.Name
	} else if len(f.name) == 0 {
		f.name = utils.Md5Str(f.url)
	}
	f.file = path.Join(config.Dir().Cache, f.name+".xml")
	// --------- 是否重新下载rss.xml ---------
	if opt.RefreshCache {
		glog.V(3).Infoln("获取Rss数据开始...")
		err := utils.HttpGet(f.url, f.file)
		if err != nil {
			glog.Errorln(err)
			return
		}
		glog.V(3).Infoln("获取Rss数据成功！")
	}
	// --------- 解析本地rss.xml ---------
	file, err := os.Open(f.file)
	if err != nil {
		glog.Errorln(err)
		return
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	if err != nil {
		glog.Errorln(err)
		return
	}
	f.feed = feed
}

// ParseBangumiAll
//  Description 第二步，获取全部番剧信息
//  Receiver f *Rss
//  Param bangumi bangumi.Bangumi 解析器
//
func (f *Rss) ParseBangumiAll(bangumi bangumi.Bangumi) []*model.Bangumi {
	bgms := make([]*model.Bangumi, len(f.feed.Items))
	wg := sync.WaitGroup{}
	for i, item := range f.feed.Items {
		wg.Add(1)
		go func(i int, item *gofeed.Item) {
			defer wg.Done()
			bgms[i] = f.ParseBangumi(item, bangumi)
		}(i, item)
	}
	wg.Wait()
	return bgms
}

// ParseBangumi
//  Description 第二步，获取番剧信息
//  Receiver f *Rss
//  Param item *gofeed.Item
//  Param bangumi bangumi.Bangumi 解析器
//
func (f *Rss) ParseBangumi(item *gofeed.Item, bangumi bangumi.Bangumi) *model.Bangumi {
	// TODO: 读取缓存，若存在且不过期则直接返回
	bgmInfo := bangumi.Parse(&model.BangumiParseOptions{
		Url:  item.Link,
		Name: item.Title,
	})
	// TODO: 写入缓存，需要线程安全
	return bgmInfo
}

package mikan

import (
	"GoBangumi/internal/core/feed"
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"GoBangumi/utils"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
	"os"
	"path"
	"regexp"
)

type Rss struct {
}

func NewRss() feed.Feed {
	return &Rss{}
}

// Parse
//  Description 第一步，解析rss
//  Receiver f *Rss
//  Param opt *models.FeedParseOptions 若RefreshCache为false，则仅重新解析本地缓存rss
//
func (f *Rss) Parse(opt *models.FeedParseOptions) []*models.FeedItem {
	if len(opt.Name) == 0 {
		opt.Name = utils.Md5Str(opt.Url)
	}
	filename := path.Join(store.Config.Setting.CachePath, opt.Name+".xml")
	// --------- 是否重新下载rss.xml ---------
	if opt.RefreshCache {
		zap.S().Info("获取Rss数据开始...")
		err := utils.HttpGet(opt.Url, filename, store.Config.Proxy())
		if err != nil {
			zap.S().Warn(err)
			return nil
		}
		zap.S().Info("获取Rss数据成功！")
	}
	// --------- 解析本地rss.xml ---------
	file, err := os.Open(filename)
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	regx := regexp.MustCompile(`<pubDate>(.*?)T`)

	var date string
	items := make([]*models.FeedItem, len(feed.Items))
	for i, item := range feed.Items {
		strs := regx.FindStringSubmatch(item.Custom["torrent"])
		if len(strs) < 2 {
			date = ""
		} else {
			date = strs[1]
		}
		_, hash := path.Split(item.Enclosures[0].URL)
		if len(hash) < 40 {
			zap.S().Warn(err)
			hash = ""
		} else {
			hash = hash[:40]
		}

		items[i] = &models.FeedItem{
			Url:     item.Link,
			Name:    item.Title,
			Date:    date,
			Torrent: item.Enclosures[0].URL,
			Hash:    hash,
		}
	}
	return items

}

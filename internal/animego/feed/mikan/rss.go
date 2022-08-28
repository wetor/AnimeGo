// Package mikan
// @Description: 获取并解析mikan rss
package mikan

import (
	"GoBangumi/internal/animego/feed"
	"GoBangumi/internal/models"
	"GoBangumi/internal/store"
	"GoBangumi/internal/utils"
	"GoBangumi/pkg/request"
	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
	"os"
	"path"
	"regexp"
)

type Rss struct {
	url  string
	name string
}

func NewRss(url, name string) feed.Feed {
	if len(name) == 0 {
		name = utils.Md5Str(url)
	}
	return &Rss{
		url:  url,
		name: name,
	}
}

// Parse
//  Description 第一步，解析rss
//  Receiver f *Rss
//
func (f *Rss) Parse() []*models.FeedItem {

	filename := path.Join(store.Config.Setting.CachePath, f.name+".xml")
	// --------- 下载rss.xml ---------
	zap.S().Info("获取Rss数据开始...")
	err := request.Get(&request.Param{
		Uri:      f.url,
		Proxy:    store.Config.Proxy(),
		SaveFile: filename,
	})
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	zap.S().Info("获取Rss数据成功！")

	// --------- 解析本地rss.xml ---------
	file, err := os.Open(filename)
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feeds, err := fp.Parse(file)
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	regx := regexp.MustCompile(`<pubDate>(.*?)T`)

	var date string
	items := make([]*models.FeedItem, len(feeds.Items))
	for i, item := range feeds.Items {
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

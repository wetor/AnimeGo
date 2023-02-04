// Package rss
// @Description: 获取并解析rss
package rss

import (
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/mmcdole/gofeed"

	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
)

type Rss struct {
	url  string
	name string
}

func NewRss(url, name string) *Rss {

	if len(name) == 0 {
		if len(url) == 0 {
			name = ""
		} else {
			name = utils.Md5Str(url)
		}

	}
	return &Rss{
		url:  url,
		name: name,
	}
}

// Parse
//  @Description: 解析rss
//  @receiver *Rss
//  @param opts ...any
//    [0] string 解析本地文件路径
//  @return items []*models.FeedItem
//
func (f *Rss) Parse(opts ...any) (items []*models.FeedItem) {
	if len(f.url) == 0 && len(opts) == 0 {
		return nil
	}
	var filename string
	log.Infof("获取Rss数据开始...")
	if len(opts) == 1 {
		filename = opts[0].(string)
	} else {
		filename = path.Join(feed.TempPath, f.name+".xml")
		// --------- 下载rss.xml ---------
		err := request.GetFile(f.url, filename)
		if err != nil {
			log.Debugf("", err)
			log.Warnf("请求Rss失败")
		}
	}
	log.Infof("获取Rss数据成功！")

	// --------- 解析本地rss.xml ---------
	file, err := os.Open(filename)
	if err != nil {
		log.Debugf("", err)
		log.Warnf("打开Rss文件失败")
	}

	defer file.Close()
	fp := gofeed.NewParser()
	feeds, err := fp.Parse(file)
	if err != nil {
		log.Debugf("", err)
		log.Warnf("解析ss失败")
	}

	regx := regexp.MustCompile(`<pubDate>(.*?)T`)
	var date string
	var length int64
	items = make([]*models.FeedItem, len(feeds.Items))
	for i, item := range feeds.Items {
		strs := regx.FindStringSubmatch(item.Custom["torrent"])
		if len(strs) < 2 {
			date = ""
		} else {
			date = strs[1]
		}

		if len(item.Enclosures) == 0 {
			log.Warnf("Torrent Enclosures错误，%s，跳过", item.Title)
			continue
		}

		length, err = strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
		if err != nil {
			log.Debugf("", err)
		}
		if length == 0 {
			log.Warnf("Torrent Length错误")
		}

		items[i] = &models.FeedItem{
			Url:      item.Link,
			Name:     item.Title,
			Date:     date,
			Type:     item.Enclosures[0].Type,
			Download: item.Enclosures[0].URL,
			Length:   length,
		}
	}
	return items

}

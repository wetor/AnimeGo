// Package rss
// @Description: 获取并解析rss
package rss

import (
	"bytes"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/mmcdole/gofeed"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
)

type Options struct {
	Url  string
	File string
	Raw  string
}

type Rss struct {
	url  string
	file string
	raw  string
}

func NewRss(opts *Options) *Rss {
	return &Rss{
		url:  opts.Url,
		file: opts.File,
		raw:  opts.Raw,
	}
}

// Parse
//
//	@Description: 解析rss
//	@receiver *Rss
//	@return items []*models.FeedItem
func (f *Rss) Parse() (items []*models.FeedItem) {
	if len(f.url) == 0 && len(f.file) == 0 && len(f.raw) == 0 {
		return nil
	}
	data := bytes.NewBuffer(nil)

	log.Infof("获取Rss数据开始...")
	if len(f.file) != 0 {
		file, err := os.Open(f.file)
		if err != nil {
			log.Debugf("", err)
			log.Warnf("打开Rss文件失败")
		}
		_, err = io.Copy(data, file)
		if err != nil {
			log.Debugf("", err)
			return nil
		}
		_ = file.Close()
	} else if len(f.raw) != 0 {
		data.WriteString(f.raw)
	} else {
		err := request.GetWriter(f.url, data)
		if err != nil {
			log.Debugf("", err)
			log.Warnf("请求Rss失败")
			return nil
		}
	}
	log.Infof("获取Rss数据成功！")

	fp := gofeed.NewParser()
	feeds, err := fp.Parse(data)
	if err != nil {
		log.Debugf("", err)
		log.Warnf("解析ss失败")
		return nil
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

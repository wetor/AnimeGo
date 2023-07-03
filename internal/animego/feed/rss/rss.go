// Package rss
// @Description: 获取并解析rss
package rss

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/internal/exceptions"
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
func (f *Rss) Parse() (items []*models.FeedItem, err error) {
	data := bytes.NewBuffer(nil)

	log.Infof("获取Rss数据开始...")
	if len(f.file) != 0 {
		file, err := os.ReadFile(f.file)
		if err != nil {
			log.DebugErr(err)
			return nil, errors.WithStack(&exceptions.ErrFeed{Message: "打开Rss文件失败"})
		}
		data.Write(file)
	} else if len(f.raw) != 0 {
		data.WriteString(f.raw)
	} else if len(f.url) != 0 {
		err := request.GetWriter(f.url, data)
		if err != nil {
			log.DebugErr(err)
			return nil, errors.WithStack(&exceptions.ErrRequest{Name: "Rss"})
		}
	} else {
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "Rss为空"})
	}
	log.Infof("获取Rss数据成功！")

	fp := gofeed.NewParser()
	feeds, err := fp.Parse(data)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "解析Rss失败"})
	}

	regx := regexp.MustCompile(`<pubDate>(.*?)T`)
	items = make([]*models.FeedItem, 0, len(feeds.Items))
	for _, item := range feeds.Items {
		if len(item.Enclosures) == 0 {
			log.Debugf("解析Rss项目「%s」详细信息失败，跳过", item.Title)
			continue
		}

		length, err := strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
		if err != nil {
			log.Debugf("解析Rss项目「%s」下载大小失败，默认0", item.Title)
		}

		var date string
		dateMatch := regx.FindStringSubmatch(item.Custom["torrent"])
		if len(dateMatch) >= 2 {
			date = dateMatch[1]
		}
		items = append(items, &models.FeedItem{
			Url:      item.Link,
			Name:     item.Title,
			Date:     date,
			Type:     item.Enclosures[0].Type,
			Download: item.Enclosures[0].URL,
			Length:   length,
		})
	}
	return items, nil
}

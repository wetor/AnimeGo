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
			return nil, errors.WithStack(&exceptions.ErrFeed{Message: "请求Rss失败"})
		}
	} else {
		return nil, err
	}
	log.Infof("获取Rss数据成功！")

	fp := gofeed.NewParser()
	feeds, err := fp.Parse(data)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "解析Rss失败"})
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
			log.Debugf("解析Rss项目 %s 详细信息失败，忽略", item.Title)
			continue
		}

		length, err = strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
		if err != nil {
			log.DebugErr(errors.Wrapf(err, "解析Rss项目 %s 下载大小失败，默认0", item.Title))
			length = 0
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
	return items, nil
}

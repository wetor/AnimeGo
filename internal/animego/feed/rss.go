// Package rss
// @Description: 获取并解析rss
package feed

import (
	"bytes"
	"os"
	"regexp"
	"strconv"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Rss struct {
}

func NewRss() *Rss {
	return &Rss{}
}

func (r *Rss) ParseFile(f string) (items []*models.FeedItem, err error) {
	if len(f) == 0 {
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "Rss为空"})
	}
	data := bytes.NewBuffer(nil)
	log.Infof("获取Rss数据开始...")
	file, err := os.ReadFile(f)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "打开Rss文件失败"})
	}
	data.Write(file)
	log.Infof("获取Rss数据成功！")
	return r.Parse(data.Bytes())
}

func (r *Rss) ParseUrl(uri string) (items []*models.FeedItem, err error) {
	if len(uri) == 0 {
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "Rss为空"})
	}
	data := bytes.NewBuffer(nil)
	log.Infof("获取Rss数据开始...")
	err = request.GetWriter(uri, data)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrRequest{Name: "Rss"})
	}
	log.Infof("获取Rss数据成功！")
	return r.Parse(data.Bytes())
}

func (r *Rss) Parse(raw []byte) (items []*models.FeedItem, err error) {
	if len(raw) == 0 {
		return nil, errors.WithStack(&exceptions.ErrFeed{Message: "Rss为空"})
	}
	data := bytes.NewBuffer(nil)
	data.Write(raw)

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
			MikanUrl:   item.Link,
			Name:       item.Title,
			Date:       date,
			Type:       item.Enclosures[0].Type,
			TorrentUrl: item.Enclosures[0].URL,
			Length:     length,
		})
	}
	return items, nil
}

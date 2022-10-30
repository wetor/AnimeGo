// Package mikan
// @Description: 获取并解析mikan rss
package mikan

import (
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"AnimeGo/internal/utils"
	"AnimeGo/pkg/errors"
	"AnimeGo/pkg/request"
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"
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
//  Description 第一步，解析rss
//  Receiver f *Rss
//
func (f *Rss) Parse() (items []*models.FeedItem, err error) {
	if len(f.url) == 0 {
		return nil, nil
	}
	filename := path.Join(store.Config.Advanced.Path.TempPath, f.name+".xml")
	// --------- 下载rss.xml ---------
	zap.S().Info("获取Rss数据开始...")
	err = request.GetFile(f.url, filename)
	if err != nil {
		zap.S().Debug(err)
		zap.S().Warn("请求Rss失败")
		return
	}
	zap.S().Info("获取Rss数据成功！")

	// --------- 解析本地rss.xml ---------
	file, err := os.Open(filename)
	if err != nil {
		err = errors.NewAniErrorD(err)
		zap.S().Debug(err)
		zap.S().Warn("打开Rss文件失败")
		return
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feeds, err := fp.Parse(file)
	if err != nil {
		err = errors.NewAniErrorD(err)
		zap.S().Debug(err)
		zap.S().Warn("解析ss失败")
		return
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
			zap.S().Debug(errors.NewAniErrorf("Torrent Enclosures错误，%s", item.Title))
			zap.S().Warn("Torrent Enclosures错误，跳过")
			continue
		}
		_, hash := path.Split(item.Enclosures[0].URL)
		if len(hash) < 40 {
			zap.S().Debug(errors.NewAniErrorf("Torrent URL错误，%s", item.Title))
			zap.S().Warn("Torrent URL错误")
			hash = ""
		} else {
			hash = hash[:40]
		}
		length, err = strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
		if err != nil {
			zap.S().Debug(errors.NewAniErrorD(err))
			zap.S().Warn("Torrent Length错误")
		}

		items[i] = &models.FeedItem{
			Url:     item.Link,
			Name:    item.Title,
			Date:    date,
			Torrent: item.Enclosures[0].URL,
			Hash:    hash,
			Length:  length,
		}
	}
	return items, nil

}

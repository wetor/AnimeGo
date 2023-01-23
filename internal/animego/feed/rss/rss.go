// Package rss
// @Description: 获取并解析rss
package rss

import (
	"os"
	"path"
	"regexp"
	"strconv"

	"github.com/mmcdole/gofeed"
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/animego/feed"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
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
//
//	Description 第一步，解析rss
//	Receiver f *Rss
func (f *Rss) Parse() (items []*models.FeedItem) {
	if len(f.url) == 0 {
		return nil
	}

	errMsg := ""
	defer errors.HandleAniError(func(err *errors.AniError) {
		zap.S().Debug(err)
		zap.S().Warn(errMsg)
	})

	filename := path.Join(feed.TempPath, f.name+".xml")
	// --------- 下载rss.xml ---------
	zap.S().Info("获取Rss数据开始...")
	errMsg = "请求Rss失败"
	err := request.GetFile(f.url, filename)
	errors.NewAniErrorD(err).TryPanic()
	zap.S().Info("获取Rss数据成功！")

	// --------- 解析本地rss.xml ---------
	errMsg = "打开Rss文件失败"
	file, err := os.Open(filename)
	errors.NewAniErrorD(err).TryPanic()

	defer file.Close()
	fp := gofeed.NewParser()
	errMsg = "解析ss失败"
	feeds, err := fp.Parse(file)
	errors.NewAniErrorD(err).TryPanic()

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
			zap.S().Warnf("Torrent Enclosures错误，%s，跳过", item.Title)
			continue
		}

		length, _ = strconv.ParseInt(item.Enclosures[0].Length, 10, 64)
		if length == 0 {
			zap.S().Warn("Torrent Length错误")
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

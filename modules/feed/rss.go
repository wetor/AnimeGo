package feed

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/utils"
	"github.com/golang/glog"
	"github.com/mmcdole/gofeed"
	"os"
	"path"
)

type Rss struct {
}

func NewRss() Feed {
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
	filename := path.Join(config.Dir().Cache, opt.Name+".xml")
	// --------- 是否重新下载rss.xml ---------
	if opt.RefreshCache {
		glog.V(3).Infoln("获取Rss数据开始...")
		err := utils.HttpGet(opt.Url, filename)
		if err != nil {
			glog.Errorln(err)
			return nil
		}
		glog.V(3).Infoln("获取Rss数据成功！")
	}
	// --------- 解析本地rss.xml ---------
	file, err := os.Open(filename)
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	defer file.Close()
	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	items := make([]*models.FeedItem, len(feed.Items))
	for i, item := range feed.Items {
		items[i] = &models.FeedItem{
			Url:  item.Link,
			Name: item.Title,
		}
	}
	return items

}

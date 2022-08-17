package mikan

import (
	"GoBangumi/internal/models"
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"fmt"
	"github.com/mmcdole/gofeed"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/GoBangumi/data/config/conf.yaml",
	})
	m.Run()
	fmt.Println("end")
}
func TestRss(t *testing.T) {
	rssConf := store.Config.RssMikan()
	rssFile := path.Join(store.Config.Setting.CachePath, rssConf.Name+".xml")

	//err := utils.HttpGet(rssConf.Url, rssFile, store.Config.Proxy())
	//if err != nil {
	//	panic(err)
	//}
	file, _ := os.Open(rssFile)
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	fmt.Println(feed.Items[0])
}

func TestRss_Parse(t *testing.T) {
	rssConf := store.Config.RssMikan()
	f := NewRss()
	items := f.Parse(&models.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name, // 文件名
		RefreshCache: true,
	})
	for _, b := range items {
		fmt.Println(b)
	}
}

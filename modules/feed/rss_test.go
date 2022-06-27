package feed

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/mmcdole/gofeed"
	"os"
	"path"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "5")
	flag.Parse()
	defer glog.Flush()

	m.Run()
	fmt.Println("end")
}
func TestRss(t *testing.T) {
	config.Init("../../data/config/conf.yaml")
	rssConf := config.RssMikan()
	rssFile := path.Join(config.Setting().CachePath, rssConf.Name+".xml")

	//err := utils.HttpGet(rssConf.Url, rssFile, config.Proxy())
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
	config.Init("../../data/config/conf.yaml")
	rssConf := config.RssMikan()
	f := NewRss()
	items := f.Parse(&models.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name, // 文件名
		RefreshCache: false,
	})
	for _, b := range items {
		fmt.Println(b)
	}
}

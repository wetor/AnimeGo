package feed

import (
	"GoBangumi/config"
	"GoBangumi/core/bangumi"
	"GoBangumi/model"
	"GoBangumi/utils"
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
	flag.Set("v", "4")
	flag.Parse()
	defer glog.Flush()

	m.Run()
	fmt.Println("end")
}
func TestRss(t *testing.T) {
	config.Init("../../data/config/conf.yaml")
	rssConf := config.Mikan()
	rssFile := path.Join(config.Dir().Cache, rssConf.Name+".xml")

	err := utils.HttpGet(rssConf.Url, rssFile)
	if err != nil {
		panic(err)
	}
	file, _ := os.Open(rssFile)
	defer file.Close()
	fp := gofeed.NewParser()
	feed, _ := fp.Parse(file)
	fmt.Println(feed)
}

func TestNewRss(t *testing.T) {
	config.Init("../../data/config/conf.yaml")
	rssConf := config.Mikan()

	feed := NewRss()
	feed.Parse(&model.FeedParseOptions{
		Url:          rssConf.Url,
		Name:         rssConf.Name,
		RefreshCache: true,
	})
	bgms := feed.ParseBangumiAll(&bangumi.Mikan{})
	fmt.Println(bgms)
	for _, b := range bgms {
		fmt.Println(b)
	}
}

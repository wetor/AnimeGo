package process

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()
	defer glog.Flush()
	config.Init("../data/config/conf.yaml")
	m.Run()
	fmt.Println("end")
}
func TestMikanProcessOne(t *testing.T) {
	p := NewMikan()
	bgms := p.ParseBangumi(&models.FeedItem{
		Url:  "https://mikanani.me/Home/Episode/6b23947f17c844570eee00177b870f16949900d4",
		Name: "[酷漫404][辉夜姬想让人告白 一超级浪漫一][09][1080P][WebRip][繁日双语][AVC AAC][MP4][字幕组招人内详]",
		Date: "2022-06-14",
	}, &bangumi.Mikan{})
	fmt.Println(bgms)

}
func TestMikanProcess(t *testing.T) {

	m := NewMikan()
	m.Run()
}

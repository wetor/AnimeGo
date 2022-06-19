package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
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
	config.Init("../../data/config/conf.yaml")
	m.Run()
	fmt.Println("end")
}
func TestNewThemoviedb(t *testing.T) {
	tmdb := NewThemoviedb()
	b := tmdb.Parse(&models.BangumiParseOptions{
		Name: "カードファイト!! ヴァンガード will+Dress",
		Date: "2022-07-04",
	})
	fmt.Println(b)
}

func TestTime(t *testing.T) {

	//fmt.Println(utils.GetTimeRangeDay("2022-04-10", 1))
}

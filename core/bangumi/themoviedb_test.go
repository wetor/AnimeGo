package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/model"
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
	b := tmdb.Parse(&model.BangumiParseOptions{
		Name: "本好きの下剋上～司書になるためには手段を選んでいられません～ 第3部",
		Date: "2022-04-01",
	})
	fmt.Println(b)
}

func TestTime(t *testing.T) {

	//fmt.Println(utils.GetTimeRangeDay("2022-04-10", 1))
}

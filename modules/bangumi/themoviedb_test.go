package bangumi

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/cache"
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	config.Init("../../data/config/conf.yaml")
	store.SetCache(cache.NewBolt())
	store.Cache.Open(config.Setting().CachePath)
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

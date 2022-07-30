package process

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
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
	config.Init("../data/config/conf.yaml")
	store.InitState = store.InitLoadConfig

	store.SetCache(cache.NewBolt())
	store.Cache.Open(config.Setting().CachePath)
	store.InitState = store.InitLoadCache

	store.InitState = store.InitConnectClient

	store.InitState = store.InitFinish
	m.Run()
	fmt.Println("end")
}
func TestMikanProcessOne(t *testing.T) {
	p := NewMikan()
	bgms := p.ParseBangumi(&models.FeedItem{
		Url:  "https://mikanani.me/Home/Episode/171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
		Name: "[夜莺家族&YYQ字幕组]New Doraemon 哆啦A梦新番[712][2022.06.25][AVC][1080P][GB_JP]",
		Date: "2022-06-26",
	}, &bangumi.Mikan{})
	fmt.Println(bgms, bgms.BangumiSeason, bgms.BangumiEp, bgms.BangumiExtra)

}
func TestMikanProcess(t *testing.T) {

	m := NewMikan()
	m.Run()
}

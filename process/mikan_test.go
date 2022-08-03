package process

import (
	"GoBangumi/models"
	"GoBangumi/modules/bangumi"
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	store.Init(nil)

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

	exit := make(chan bool)
	m.Run(exit)

	go func() {
		time.Sleep(10 * time.Minute)
		m.Exit()
	}()

	<-exit
}

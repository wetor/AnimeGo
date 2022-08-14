package process

import (
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
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

//func TestMikanProcessOne(t *testing.T) {
//	p := NewMikan()
//	animes := p.ParseBangumi(&models.FeedItem{
//		Url:  "https://mikanani.me/Home/Episode/171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
//		Name: "[夜莺家族&YYQ字幕组]New Doraemon 哆啦A梦新番[712][2022.06.25][AVC][1080P][GB_JP]",
//		Date: "2022-06-26",
//	}, &mikan.Mikan{})
//	fmt.Println(animes, animes.AnimeSeason, animes.AnimeEp, animes.AnimeExtra)
//
//}
func TestMikanProcess(t *testing.T) {

	m := NewMikan()
	ctx, cancel := context.WithCancel(context.Background())
	m.Run(ctx)
	store.WG.Add(2)
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()

	store.WG.Wait()
}

func TestG(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		time.Sleep(3 * time.Second)
		panic("hellp")

	}()
	wg.Wait()
}

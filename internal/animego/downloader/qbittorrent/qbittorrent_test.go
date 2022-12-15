package qbittorrent

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/test"
	"testing"
	"time"
)

var qbt *QBittorrent

func TestMain(m *testing.M) {
	fmt.Println("begin")
	test.TestInit()
	conf := store.Config.Setting.Client.QBittorrent
	qbt = NewQBittorrent(conf.Url, conf.Username, conf.Password)
	qbt.Start(context.Background())
	m.Run()
	fmt.Println("end")
}

func TestQBittorrent_Run(t *testing.T) {
	time.Sleep(3 * time.Second)
	go func() {
		for {
			fmt.Println("ver", qbt.Version())
			time.Sleep(1 * time.Second)
		}
	}()

	store.WG.Wait()
}

func Test_QBittorrent(t *testing.T) {
	time.Sleep(2 * time.Second)
	list := qbt.List(&models.ClientListOptions{
		Category: "AnimeGo",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}

}

func TestQBittorrent_Add(t *testing.T) {
	qbt.Add(&models.ClientAddOptions{
		Urls: []string{
			"https://mikanani.me/Download/20220612/4407d51f30f6033513cbe56cae0120881b0a7406.torrent",
			"https://mikanani.me/Download/20220611/56e13c0c4788b77782722ee46d3c6f27233f676b.torrent",
		},
		SavePath:    "/srv/dev-disk-by-uuid-317ec4d4-2933-4fa6-9d7b-fcdfc339de04/share/downloads/complete/",
		Category:    "test",
		Tag:         "test_tag",
		SeedingTime: 60,
	})

}

func TestQBittorrent_Delete(t *testing.T) {

	qbt.Delete(&models.ClientDeleteOptions{
		Hash: []string{"4407d51f30f6033513cbe56cae0120881b0a7406"},
	})

}

func TestQBittorrent_Get(t *testing.T) {

	a := qbt.Get(&models.ClientGetOptions{Hash: "171f3b402fa4cf770ef267c0744a81b6b9ad77f2"})
	fmt.Println(a)
}

func TestQBittorrent_Rename(t *testing.T) {
	qbt.Rename(&models.ClientRenameOptions{
		Hash:    "171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
		OldPath: "[夜莺家族&YYQ字幕组]New Doraemon 哆啦A梦新番[712][2022.06.25][AVC][1080P][GB_JP]/[夜莺家族&YYQ字幕组]New Doraemon 哆啦A梦新番[712][2022.06.25][AVC][1080P][GB_JP].mp4",
		NewPath: "[夜莺家族&YYQ字幕组]New Doraemon 哆啦A梦新番[712][2022.06.25][AVC][1080P][GB_JP]/哆啦A梦[第1季][第712集].mp4",
	})
}
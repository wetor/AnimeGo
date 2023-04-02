package qbittorrent_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

var qbt *qbittorrent.QBittorrent

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	wg := sync.WaitGroup{}
	downloader.Init(&downloader.Options{
		WG: &wg,
	})
	qbt = qbittorrent.NewQBittorrent("http://127.0.0.1:8080", "admin", "adminadmin")
	qbt.Start(context.Background())
	for i := 0; i < 5 && !qbt.Connected(); i++ {
		time.Sleep(time.Second)
	}
	m.Run()
	fmt.Println("end")
}

func Test_QBittorrent(t *testing.T) {
	t.Skip("跳过Qbittorrent测试")
	list := qbt.List(&models.ClientListOptions{
		Category: "AnimeGo",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}

}

func TestQBittorrent_Add(t *testing.T) {
	t.Skip("跳过Qbittorrent测试")
	qbt.Add(&models.ClientAddOptions{
		Urls: []string{
			"https://mikanani.me/Download/20220612/4407d51f30f6033513cbe56cae0120881b0a7406.torrent",
			"https://mikanani.me/Download/20220611/56e13c0c4788b77782722ee46d3c6f27233f676b.torrent",
		},
		SavePath:    "/tmp/test",
		Category:    "test",
		Tag:         "test_tag",
		SeedingTime: 60,
	})
	time.Sleep(3 * time.Second)
	list := qbt.List(&models.ClientListOptions{
		Category: "test",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
	qbt.Delete(&models.ClientDeleteOptions{
		Hash: []string{"4407d51f30f6033513cbe56cae0120881b0a7406", "56e13c0c4788b77782722ee46d3c6f27233f676b"},
	})

}

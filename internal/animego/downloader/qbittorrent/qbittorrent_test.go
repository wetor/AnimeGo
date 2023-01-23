package qbittorrent

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/models"
)

var qbt *QBittorrent

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	wg := sync.WaitGroup{}
	downloader.Init(&downloader.Options{
		WG: &wg,
	})
	qbt = NewQBittorrent("http://192.168.10.50:8080", "admin", "adminadmin")
	qbt.Start(context.Background())
	for !qbt.Connected() {
		time.Sleep(time.Second)
	}
	m.Run()
	fmt.Println("end")
}

func Test_QBittorrent(t *testing.T) {
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

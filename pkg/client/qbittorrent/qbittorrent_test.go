package qbittorrent_test

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/pkg/client"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/pkg/client/qbittorrent"
	"github.com/wetor/AnimeGo/pkg/log"
)

var qbt *qbittorrent.QBittorrent

func TestMain(m *testing.M) {
	fmt.Println("跳过Qbittorrent测试")
	os.Exit(0)
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	wg := sync.WaitGroup{}

	qbt = qbittorrent.NewQBittorrent(&qbittorrent.Options{
		Url:          "http://127.0.0.1:8080",
		Username:     "admin",
		Password:     "adminadmin",
		DownloadPath: "/tmp/test",
		WG:           &wg,
	})
	qbt.Start(context.Background())
	for i := 0; i < 5 && !qbt.Connected(); i++ {
		time.Sleep(time.Second)
	}
	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func Test_QBittorrent(t *testing.T) {
	t.Skip("跳过Qbittorrent测试")
	list, _ := qbt.List(&client.ListOptions{
		Category: "AnimeGo",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}

}

func TestQBittorrent_Add(t *testing.T) {
	t.Skip("跳过Qbittorrent测试")
	qbt.Add(&client.AddOptions{
		Url:         "https://mikanani.me/Download/20220612/4407d51f30f6033513cbe56cae0120881b0a7406.torrent",
		SavePath:    "/tmp/test",
		Category:    "test",
		Tag:         "test_tag",
		SeedingTime: 60,
	})
	time.Sleep(3 * time.Second)
	list, _ := qbt.List(&client.ListOptions{
		Category: "test",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
	qbt.Delete(&client.DeleteOptions{
		Hash: []string{"4407d51f30f6033513cbe56cae0120881b0a7406", "56e13c0c4788b77782722ee46d3c6f27233f676b"},
	})

}

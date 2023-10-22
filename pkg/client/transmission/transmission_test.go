package transmission_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/client/transmission"
	"github.com/wetor/AnimeGo/pkg/log"
)

var (
	tbt *transmission.Transmission
	wg  sync.WaitGroup
)

func TestMain(m *testing.M) {
	//fmt.Println("跳过Qbittorrent测试")
	//os.Exit(0)
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	wg = sync.WaitGroup{}

	tbt = transmission.NewTransmission(&transmission.Options{
		Url:          "http://127.0.0.1:9091",
		Username:     "admin",
		Password:     "adminadmin",
		DownloadPath: "C:/Users/wetor/GolandProjects/AnimeGo/download/incomplete",
		WG:           &wg,
		Ctx:          context.Background(),
	})
	tbt.Start()
	for i := 0; i < 5 && !tbt.Connected(); i++ {
		time.Sleep(time.Second)
	}
	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestTransmission_Add(t *testing.T) {
	err := tbt.Add(&client.AddOptions{
		File:        "C:/Users/wetor/Downloads/ebaa0327cd552d939485e989a1db7b11c0d38290.torrent",
		SavePath:    "C:/Users/wetor/GolandProjects/AnimeGo/download/incomplete",
		Category:    "test",
		Tag:         "test_tag",
		SeedingTime: 60,
	})
	if err != nil {
		panic(err)
	}
	list, err := tbt.List(&client.ListOptions{
		Category: "test",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
}

func TestTransmission_Start(t *testing.T) {
	wg.Wait()
}

func TestTransmission_List(t *testing.T) {

	for j := 0; j < 100; j++ {
		list, err := tbt.List(&client.ListOptions{
			Category: "test",
		})
		if err != nil {
			panic(err)
		}
		for _, i := range list {
			fmt.Println(j, i.Name, i.Hash, i.State)
		}
		time.Sleep(1 * time.Second)
	}

}

func TestTransmission_Delete(t *testing.T) {
	err := tbt.Delete(&client.DeleteOptions{
		Hash:       []string{"ebaa0327cd552d939485e989a1db7b11c0d38290"},
		DeleteFile: true,
	})
	if err != nil {
		panic(err)
	}

}

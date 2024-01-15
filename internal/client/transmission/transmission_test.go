package transmission_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/client/transmission"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "client"

var (
	tbt *transmission.Transmission
	wg  sync.WaitGroup
)

func TestMain(m *testing.M) {
	fmt.Println("跳过Qbittorrent测试")
	os.Exit(0)
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	wg = sync.WaitGroup{}

	tbt = transmission.NewTransmission(&models.ClientOptions{
		Url:      "http://127.0.0.1:9091",
		Username: "admin",
		Password: "adminadmin",

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

func TestTransmission_AddFile(t *testing.T) {
	err := tbt.Add(&models.AddOptions{
		File:     test.GetDataPath(testdata, "25a0ee3f251d052507df7a552673711bb79279c6.torrent"),
		SavePath: "C:/Users/wetor/GolandProjects/AnimeGo/download/incomplete",
		Category: "test",
		Tag:      "test_tag",
	})
	if err != nil {
		panic(err)
	}

	err = tbt.Add(&models.AddOptions{
		File:     test.GetDataPath(testdata, "f690ffb1419152efc81b332dafb3456f9ece1744.torrent"),
		SavePath: "C:/Users/wetor/GolandProjects/AnimeGo/download/incomplete",
		Category: "test2",
		Tag:      "test_tag2",
	})
	if err != nil {
		panic(err)
	}
	// 25a0ee3f251d052507df7a552673711bb79279c6
	list, _ := tbt.List(&models.ListOptions{
		Category: "test",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
	// 25a0ee3f251d052507df7a552673711bb79279c6
	list, _ = tbt.List(&models.ListOptions{
		Tag:      "test_tag",
		Category: "test",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
	// f690ffb1419152efc81b332dafb3456f9ece1744
	list, _ = tbt.List(&models.ListOptions{
		Category: "test2",
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
	// f690ffb1419152efc81b332dafb3456f9ece1744
	list, _ = tbt.List(&models.ListOptions{
		Tag:      "test_tag2",
		Category: "test2",
		Status:   transmission.TorrentStatusDownload,
	})
	fmt.Println(len(list))
	for _, i := range list {
		fmt.Println(i.Name, i.Hash, i.State)
	}
}

func TestTransmission_AddUrl(t *testing.T) {
	err := tbt.Add(&models.AddOptions{
		Url:      "magnet:?xt=urn:btih:f690ffb1419152efc81b332dafb3456f9ece1744&tr=http%3a%2f%2ft.nyaatracker.com%2fannounce&tr=http%3a%2f%2ftracker.kamigami.org%3a2710%2fannounce&tr=http%3a%2f%2fshare.camoe.cn%3a8080%2fannounce&tr=http%3a%2f%2fopentracker.acgnx.se%2fannounce&tr=http%3a%2f%2fanidex.moe%3a6969%2fannounce&tr=http%3a%2f%2ft.acg.rip%3a6699%2fannounce&tr=https%3a%2f%2ftr.bangumi.moe%3a9696%2fannounce&tr=udp%3a%2f%2ftr.bangumi.moe%3a6969%2fannounce&tr=http%3a%2f%2fopen.acgtracker.com%3a1096%2fannounce&tr=udp%3a%2f%2ftracker.opentrackr.org%3a1337%2fannounce",
		SavePath: "C:/Users/wetor/GolandProjects/AnimeGo/download/incomplete",
		Category: "test3",
		Tag:      "test_tag3",
	})
	if err != nil {
		panic(err)
	}
	list, err := tbt.List(&models.ListOptions{
		Category: "test3",
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
	categories := []string{"test", "test2"}
	for _, c := range categories {
		list, err := tbt.List(&models.ListOptions{
			Category: c,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(c)
		for _, i := range list {
			fmt.Println(i)
		}
		time.Sleep(1 * time.Second)
	}

}

func TestTransmission_Delete(t *testing.T) {
	err := tbt.Delete(&models.DeleteOptions{
		Hash:       []string{"25a0ee3f251d052507df7a552673711bb79279c6"},
		DeleteFile: true,
	})
	if err != nil {
		panic(err)
	}

}

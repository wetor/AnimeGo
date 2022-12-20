package downloader

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/test"
	"path"
	"testing"
	"time"
)

var qbt downloader.Client

func TestMain(m *testing.M) {
	fmt.Println("begin")
	test.TestInit()

	conf := store.Config.Setting.Client.QBittorrent
	qbt = qbittorrent.NewQBittorrent(conf.Url, conf.Username, conf.Password)
	qbt.Start(context.Background())

	m.Run()
	fmt.Println("end")
}

func Download1(m *Manager) {
	animes := &models.AnimeEntity{
		ID:      18692,
		Name:    "ドラえもん",
		NameCN:  "哆啦A梦",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep:      712,
		Date:    "2022-06-27",
		EpID:    1114366,
		MikanID: 681,
		DownloadInfo: &models.DownloadInfo{
			Url:  "https://mikanani.me/Download/20220626/171f3b402fa4cf770ef267c0744a81b6b9ad77f2.torrent",
			Hash: "171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
		},
	}
	m.Download(animes)
}

func Download2(m *Manager) {

	animes := &models.AnimeEntity{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "海贼王",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep:      1026,
		Date:    "2022-06-27",
		EpID:    1114366,
		MikanID: 228,
		DownloadInfo: &models.DownloadInfo{
			Url:  "https://mikanani.me/Download/20220725/193f881098f1a2a4347e8b04512118090f79345d.torrent",
			Hash: "193f881098f1a2a4347e8b04512118090f79345d",
		},
	}
	m.Download(animes)
}

func TestManager_Update(t *testing.T) {
	m := NewManager(qbt, store.Cache, nil)
	Download1(m)
	Download2(m)
	exit := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)

	go func() {
		time.Sleep(30 * time.Second)
		cancel()
		exit <- true
	}()

	time.Sleep(10 * time.Second)

	<-exit

}

func TestNewManager(t *testing.T) {
	urls := "https://mikanani.me/Download/20220701/2f71409a4535f23e15b1cfe054342589ca951a68.torrent"
	_, file := path.Split(urls)
	hash := file[:40]
	fmt.Println(hash, len(hash))
}

func TestManager2_Update(t *testing.T) {
	m := NewManager(qbt, store.Cache, nil)
	//animes := &models.AnimeEntity{
	//	ID:      18692,
	//	Name:    "ドラえもん",
	//	NameCN:  "哆啦A梦444",
	//	AirDate: "2005-04-15",
	//	Eps:     0,
	//	Season:  1,
	//	Ep:      712,
	//	Date:    "2022-06-27",
	//	EpID:    1114366,
	//	MikanID: 681,
	//	DownloadInfo: &models.DownloadInfo{
	//		Url:  "https://mikanani.me/Download/20221126/1474455a9238e328691c4eadbbab8cde7191df7a.torrent",
	//		Hash: "1474455a9238e328691c4eadbbab8cde7191df7a",
	//	},
	//}
	//m.Download(animes)
	exit := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)

	go func() {
		time.Sleep(200 * time.Second)
		cancel()
		exit <- true
	}()

	time.Sleep(10 * time.Second)

	<-exit

}

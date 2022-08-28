package downloader

import (
	"AnimeGo/internal/animego/downloader"
	"AnimeGo/internal/animego/downloader/qbittorent"
	"AnimeGo/internal/logger"
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"context"
	"fmt"
	"path"
	"testing"
	"time"
)

var qbt downloader.Client

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()

	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/AnimeGo/data/config/conf.yaml",
	})

	conf := store.Config.ClientQBt()
	qbt = qbittorent.NewQBittorrent(conf.Url, conf.Username, conf.Password, nil)

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
		AnimeSeason: &models.AnimeSeason{
			Season: 1,
		},
		AnimeEp: &models.AnimeEp{
			Ep:   712,
			Date: "2022-06-27",
			EpID: 1114366,
		},
		AnimeExtra: &models.AnimeExtra{
			MikanID:  681,
			MikanUrl: "https://mikanani.me/Home/Episode/171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
		},
		TorrentInfo: &models.TorrentInfo{
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
		AnimeSeason: &models.AnimeSeason{
			Season: 1,
		},
		AnimeEp: &models.AnimeEp{
			Ep:   1026,
			Date: "2022-06-27",
			EpID: 1114366,
		},
		AnimeExtra: &models.AnimeExtra{
			MikanID:  228,
			MikanUrl: "https://mikanani.me/Home/Episode/193f881098f1a2a4347e8b04512118090f79345d",
		},
		TorrentInfo: &models.TorrentInfo{
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

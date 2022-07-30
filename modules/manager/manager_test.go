package manager

import (
	"GoBangumi/config"
	"GoBangumi/models"
	"GoBangumi/modules/cache"
	"GoBangumi/modules/client"
	"GoBangumi/store"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"path"
	"testing"
	"time"
)

var qbt client.Client

func TestMain(m *testing.M) {
	fmt.Println("begin")
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "log")
	flag.Set("v", "10")
	flag.Parse()
	defer glog.Flush()

	config.Init("../../data/config/conf.yaml")
	store.InitState = store.InitLoadConfig

	store.SetCache(cache.NewBolt())
	store.Cache.Open(config.Setting().CachePath)
	store.InitState = store.InitLoadCache

	conf := config.ClientQBt()
	qbt = client.NewQBittorrent(conf.Url, conf.Username, conf.Password)

	store.InitState = store.InitConnectClient

	store.InitState = store.InitFinish
	m.Run()
	fmt.Println("end")
}
func Download1(m *Manager) {
	bgms := &models.Bangumi{
		ID:      18692,
		Name:    "ドラえもん",
		NameCN:  "哆啦A梦",
		AirDate: "2005-04-15",
		Eps:     0,
		BangumiSeason: &models.BangumiSeason{
			Season: 1,
		},
		BangumiEp: &models.BangumiEp{
			Ep:   712,
			Date: "2022-06-27",
			EpID: 1114366,
		},
		BangumiExtra: &models.BangumiExtra{
			SubID:  681,
			SubUrl: "https://mikanani.me/Home/Episode/171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
		},
		TorrentInfo: &models.TorrentInfo{
			Url:  "https://mikanani.me/Download/20220626/171f3b402fa4cf770ef267c0744a81b6b9ad77f2.torrent",
			Hash: "171f3b402fa4cf770ef267c0744a81b6b9ad77f2",
		},
	}
	m.Download(bgms)
}

func Download2(m *Manager) {

	bgms := &models.Bangumi{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "海贼王",
		AirDate: "2005-04-15",
		Eps:     0,
		BangumiSeason: &models.BangumiSeason{
			Season: 1,
		},
		BangumiEp: &models.BangumiEp{
			Ep:   1026,
			Date: "2022-06-27",
			EpID: 1114366,
		},
		BangumiExtra: &models.BangumiExtra{
			SubID:  228,
			SubUrl: "https://mikanani.me/Home/Episode/193f881098f1a2a4347e8b04512118090f79345d",
		},
		TorrentInfo: &models.TorrentInfo{
			Url:  "https://mikanani.me/Download/20220725/193f881098f1a2a4347e8b04512118090f79345d.torrent",
			Hash: "193f881098f1a2a4347e8b04512118090f79345d",
		},
	}
	m.Download(bgms)
}

func TestManager_Update(t *testing.T) {
	m := NewManager(qbt)
	Download1(m)
	Download2(m)
	exit := make(chan bool)
	m.Start(exit)

	go func() {
		time.Sleep(30 * time.Second)
		m.Exit()
	}()

	time.Sleep(10 * time.Second)

	<-exit

}

func TestManager_Get(t *testing.T) {

}

func TestNewManager(t *testing.T) {
	urls := "https://mikanani.me/Download/20220701/2f71409a4535f23e15b1cfe054342589ca951a68.torrent"
	_, file := path.Split(urls)
	hash := file[:40]
	fmt.Println(hash, len(hash))
}

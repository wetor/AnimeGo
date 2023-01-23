package downloader

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/downloader/qbittorrent"
	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/cache"
)

var (
	qbt downloader.Client
	mgr *Manager
	wg  sync.WaitGroup

	ctx, cancel = context.WithCancel(context.Background())
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	_ = utils.CreateMutiDir("data")

	downloader.Init(&downloader.Options{
		WG: &wg,
	})
	qbt = qbittorrent.NewQBittorrent("http://192.168.10.50:8080", "admin", "adminadmin")
	qbt.Start(ctx)

	manager.Init(&manager.Options{
		Downloader: manager.Downloader{
			UpdateDelaySecond:      5,
			DownloadPath:           "/tmp/download",
			SavePath:               "/tmp/save",
			Category:               "test",
			Tag:                    "test",
			AllowDuplicateDownload: false,
			SeedingTimeMinute:      0,
			IgnoreSizeMaxKb:        100,
			Rename:                 "wait_move",
		},
		WG: &wg,
	})
	b := cache.NewBolt()
	b.Open("data/test.db")
	mgr = NewManager(qbt, b, nil)

	for !qbt.Connected() {
		time.Sleep(time.Second)
	}
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
		MikanID: 228,
		DownloadInfo: &models.DownloadInfo{
			Url:  "https://mikanani.me/Download/20220725/193f881098f1a2a4347e8b04512118090f79345d.torrent",
			Hash: "193f881098f1a2a4347e8b04512118090f79345d",
		},
	}
	m.Download(animes)
}

func TestManager_Update(t *testing.T) {
	Download1(mgr)
	Download2(mgr)

	mgr.Start(ctx)

	go func() {
		time.Sleep(30 * time.Second)
		mgr.client.Delete(&models.ClientDeleteOptions{
			Hash:       []string{"171f3b402fa4cf770ef267c0744a81b6b9ad77f2", "193f881098f1a2a4347e8b04512118090f79345d"},
			DeleteFile: true,
		})
		cancel()
		wg.Done()
		os.Remove("data/test.db")
	}()
	wg.Wait()

}

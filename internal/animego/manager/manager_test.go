package manager_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/animego/rename"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
	ContentFile  = "file.mp4"
)

var (
	qbt         api.Downloader
	qbtConnect  = true
	mgr         *manager.Manager
	wg          sync.WaitGroup
	ctx, cancel = context.WithCancel(context.Background())
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	_ = utils.CreateMutiDir(DownloadPath)
	_ = utils.CreateMutiDir(SavePath)
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	qbt = &ClientMock{}
	qbt.Start(ctx)

	manager.Init(&manager.Options{
		Downloader: manager.Downloader{
			UpdateDelaySecond:      1,
			DownloadPath:           DownloadPath,
			SavePath:               SavePath,
			Category:               "test",
			Tag:                    "test",
			AllowDuplicateDownload: false,
			SeedingTimeMinute:      0,
			IgnoreSizeMaxKb:        1,
			Rename:                 "wait_move",
		},
		WG: &wg,
	})
	b := cache.NewBolt()
	b.Open("data/test.db")
	b.Add("name2status")
	b.Put("name2status", "test", &models.DownloadStatus{
		Hash:       "0000a4042b0bac2406b71023fdfe5e9054ebb832",
		State:      "complete",
		Path:       SavePath + "/test/test.mp4",
		Init:       true,
		Renamed:    true,
		Downloaded: true,
		Scraped:    true,
		Seeded:     true,
	}, 0)

	mgr = manager.NewManager(qbt, b, &rename.Rename{}, nil)

	mgr.Start(ctx)
	m.Run()

	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func Download1(m *manager.Manager) *models.AnimeEntity {
	tempHash := "4199a4042b0bac2406b71023fdfe5e9054ebb832"
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ドラえもん",
		NameCN:  "动画1",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep:      712,
		MikanID: 681,
		DownloadInfo: &models.DownloadInfo{
			Hash: tempHash,
		},
	}
	name2hash[anime.FullName()] = tempHash
	m.Download(anime)
	return anime
}

func Download2(m *manager.Manager) *models.AnimeEntity {
	tempHash := "6666a4042b0bac2406b71023fdfe5e9054ebb832"
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "动画2",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep:      1026,
		MikanID: 228,
		DownloadInfo: &models.DownloadInfo{
			Hash: tempHash,
		},
	}
	name2hash[anime.FullName()] = tempHash
	m.Download(anime)
	return anime
}

func Download3(m *manager.Manager) *models.AnimeEntity {
	tempHash := "7777a4042b0bac2406b71023fdfe5e9054ebb832"
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "动画3",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep:      996,
		MikanID: 228,
		DownloadInfo: &models.DownloadInfo{
			Hash: tempHash,
		},
	}
	name2hash[anime.FullName()] = tempHash
	m.Download(anime)
	return anime
}

func TestManager_Start(t *testing.T) {
	go func() {
		time.Sleep(15 * time.Second)
		cancel()
		time.Sleep(1*time.Second + 500*time.Millisecond)
		wg.Done()
	}()
	fmt.Println("下载 1")
	a1 := Download1(mgr)
	fmt.Println("下载 2")
	a2 := Download2(mgr)
	time.Sleep(2*time.Second + 500*time.Millisecond)
	fmt.Println("删除 2")
	mgr.DeleteCache(a2.FullName())
	qbtConnect = false
	fmt.Println("下载 3")
	a3 := Download3(mgr)
	time.Sleep(2*time.Second + 500*time.Millisecond)
	manager.Conf.Rename = "link_delete"
	qbtConnect = true
	time.Sleep(1*time.Second + 500*time.Millisecond)
	fmt.Println("重复下载 1")
	Download1(mgr)
	fmt.Println("重复下载 3")
	Download3(mgr)

	wg.Wait()
	assert.FileExists(t, xpath.Join(SavePath, a1.DirName(), a1.FileName()+xpath.Ext(ContentFile)))
	assert.FileExists(t, xpath.Join(SavePath, a3.DirName(), a3.FileName()+xpath.Ext(ContentFile)))
}

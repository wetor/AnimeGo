package downloader_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/manager"
	downloaderMgr "github.com/wetor/AnimeGo/internal/animego/manager/downloader"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
	ContentFile  = "file.mp4"
)

var (
	qbt         api.Downloader
	mgr         *downloaderMgr.Manager
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
	mgr = downloaderMgr.NewManager(qbt, b, nil)

	m.Run()

	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func Download1(m *downloaderMgr.Manager) *models.AnimeEntity {
	tempHash := utils.RandString(40)
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ドラえもん",
		NameCN:  "哆啦A梦",
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

func Download2(m *downloaderMgr.Manager) *models.AnimeEntity {
	tempHash := utils.RandString(40)
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "海贼王",
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

func TestManager_Start(t *testing.T) {
	mgr.Start(ctx)
	go func() {
		time.Sleep(15 * time.Second)
		cancel()
		wg.Done()
	}()
	a1 := Download1(mgr)
	time.Sleep(2 * time.Second)
	a2 := Download2(mgr)
	wg.Wait()
	assert.FileExists(t, xpath.Join(SavePath, a1.DirName(), a1.FileName()+xpath.Ext(ContentFile)))
	assert.FileExists(t, xpath.Join(SavePath, a2.DirName(), a2.FileName()+xpath.Ext(ContentFile)))
}

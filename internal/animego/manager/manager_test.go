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
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/torrent"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
	ContentFile  = "file.mp4"
)

var (
	qbt        *ClientMock
	qbtConnect = true
	rename     *renamer.Manager
	mgr        *manager.Manager
	db         *cache.Bolt
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
	plugin.Init(&plugin.Options{
		Path:  "../../../assets/plugin",
		Debug: true,
	})
	qbt = &ClientMock{}
	wg := sync.WaitGroup{}
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
	db = cache.NewBolt()
	db.Open("data/test.db")
	db.Add("name2status")
	db.Put("name2status", "test[第1季][第1集]", &models.DownloadStatus{
		Hash:       "0000a4042b0bac2406b71023fdfe5e9054ebb832",
		State:      "complete",
		Path:       SavePath + "/test/test.mp4",
		Init:       true,
		Renamed:    true,
		Downloaded: true,
		Scraped:    true,
		Seeded:     true,
		ExpireAt:   0,
	}, 0)
	renamer.Init(&renamer.Options{
		WG:                &wg,
		UpdateDelaySecond: 1,
	})

	rename = renamer.NewManager(renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "rename/builtin_rename.py",
	}))

	mgr = manager.NewManager(qbt, db, rename)

	m.Run()

	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func download(name string, season, ep int) (file, fullname, hash string) {
	hash = utils.MD5([]byte(fmt.Sprintf("%v%v%v", name, season, ep)))
	anime := &models.AnimeEntity{
		NameCN: name,
		Season: season,
		Ep:     ep,
		Torrent: &torrent.Torrent{
			Hash: hash,
		},
	}
	fullname = anime.FullName()
	qbt.AddName(fullname, hash)
	mgr.Download(anime)
	file = xpath.Join(SavePath, anime.DirName(), anime.FileName()+xpath.Ext(ContentFile))
	return
}

func initTest() (*sync.WaitGroup, func()) {
	wg := sync.WaitGroup{}
	renamer.WG = &wg
	manager.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())
	qbt.Init()
	qbt.Start(ctx)
	rename.Start(ctx)
	mgr.Start(ctx)
	return &wg, cancel
}

func TestManager_Success(t *testing.T) {
	wg, cancel := initTest()

	var file1 string
	go func() {
		time.Sleep(7 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, 1)
	}
	wg.Wait()
	assert.FileExists(t, file1)
	_ = os.Remove(file1)
}

func TestManager_Exist(t *testing.T) {
	wg, cancel := initTest()

	var file1 string
	go func() {
		time.Sleep(7 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, 1)

	}
	{
		log.Info("下载 1")
		file1, _, _ = download("动画1", 1, 1)
	}
	wg.Wait()
	assert.FileExists(t, file1)
	_ = os.Remove(file1)
}

func TestManager_DeleteFile_ReDownload(t *testing.T) {
	wg, cancel := initTest()

	var file1 string
	go func() {
		time.Sleep(13 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, 1)

		}
		time.Sleep(6*time.Second + 300*time.Millisecond)
		{
			log.Info("删除 1 文件")
			_ = os.Remove(file1)
			assert.NoFileExists(t, file1)
		}
		time.Sleep(500 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, 1)
		}
	}()
	wg.Wait()
	assert.FileExists(t, file1)
	_ = os.Remove(file1)
}

func TestManager_DeleteCache_ReDownload(t *testing.T) {
	wg, cancel := initTest()

	var file1 string
	var name1 string
	go func() {
		time.Sleep(13 * time.Second)
		cancel()
	}()
	manager.Conf.Rename = "move"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, name1, _ = download("动画1", 1, 1)

		}
		time.Sleep(6*time.Second + 300*time.Millisecond)
		{
			log.Info("删除 1 缓存")
			mgr.DeleteCache(name1)
			assert.FileExists(t, file1)
		}
		time.Sleep(500 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, 1)
		}
	}()
	wg.Wait()
	assert.FileExists(t, file1)
	_ = os.Remove(file1)
}

func TestManager_WaitClient(t *testing.T) {
	wg, cancel := initTest()

	var file1, file2 string
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	go func() {
		{
			log.Info("Client离线")
			qbtConnect = false
		}
		time.Sleep(2 * time.Second)
		{
			log.Info("Client恢复")
			qbtConnect = true
		}
	}()
	manager.Conf.Rename = "link_delete"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, 1)

		}
		time.Sleep(1*time.Second + 300*time.Millisecond)
		{
			log.Info("下载 2")
			file2, _, _ = download("动画2", 1, 1)
		}
	}()

	wg.Wait()
	assert.FileExists(t, file1)
	assert.FileExists(t, file2)
	_ = os.Remove(file1)
	_ = os.Remove(file2)
}

func TestManager_WaitClient_FullChan(t *testing.T) {
	wg, cancel := initTest()

	var file []string
	go func() {
		time.Sleep(12 * time.Second)
		cancel()
	}()
	go func() {
		{
			log.Info("Client离线")
			qbtConnect = false
		}
		time.Sleep(5 * time.Second)
		{
			log.Info("Client恢复")
			qbtConnect = true
		}
	}()
	manager.Conf.Rename = "move"
	go func() {
		time.Sleep(300 * time.Millisecond)
		{
			for i := 1; i <= manager.DownloadChanDefaultCap; i++ {
				log.Infof("下载 %d", i)
				f, _, _ := download("动画1", 1, i)
				file = append(file, f)
			}
		}
		{
			log.Info("下载 2")
			f, _, _ := download("动画2", 1, 1)
			file = append(file, f)
			log.Info("下载 3")
			f, _, _ = download("动画3", 1, 1)
			file = append(file, f)
		}
	}()

	wg.Wait()
	for _, f := range file {
		assert.FileExists(t, f)
		_ = os.Remove(f)
	}
}

func TestManager_ReStart_NotDownloaded(t *testing.T) {
	var file1 string
	{
		wg, cancel := initTest()
		go func() {
			time.Sleep(5 * time.Second)
			cancel()
		}()
		manager.Conf.Rename = "move"
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, 1)

		}
		wg.Wait()
		assert.FileExists(t, file1)
	}
	time.Sleep(1 * time.Second)
	{
		log.Info("重启")
		wg, cancel := initTest()
		go func() {
			time.Sleep(3 * time.Second)
			cancel()
		}()
		{
			log.Info("下载 1")
			file1, _, _ = download("动画1", 1, 1)
		}
		wg.Wait()
	}

	assert.FileExists(t, file1)
	_ = os.Remove(file1)

}

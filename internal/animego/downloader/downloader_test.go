package downloader_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/clientnotifier"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	"github.com/wetor/AnimeGo/internal/client/qbittorrent"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/wire"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"

	HookTimeUnix = 100
)

var (
	qbt          *qbittorrent.ClientMock
	dbs          *clientnotifier.Notifier
	mgr          *downloader.Manager
	rename       *renamer.Manager
	renamePlugin *renamer.Rename
	db           *cache.Bolt
	out          *bytes.Buffer
)

var defaultUpdateList = func(m *qbittorrent.ClientMock) {
	for name, item := range m.Name2item {
		switch item.State {
		case qbittorrent.QbtDownloading, qbittorrent.QbtQueuedUP:
			item.State = qbittorrent.QbtDownloading
			item.Progress += 0.2
			if item.Progress >= 1.005 {
				item.Progress = 0
				item.State = qbittorrent.QbtUploading
			} else {
				log.Debugf("%s: 下载: %v (%s)", name, int((item.Progress+0.005)*100), item.State)
			}
		case qbittorrent.QbtUploading:
			item.Progress += 0.2
			if item.Progress >= 1.005 {
				item.State = qbittorrent.QbtCheckingUP
			} else {
				log.Debugf("%s: 做种: %v (%s)", name, int((item.Progress+0.005)*100), item.State)
			}
		default:
			log.Debugf("%s: %s", name, item.State)
		}
	}
}

var waitUpdateList = func(m *qbittorrent.ClientMock) {
	for name, item := range m.Name2item {
		switch item.State {
		case qbittorrent.QbtQueuedUP:
			item.Progress += 0.4
			if item.Progress >= 1.005 {
				item.Progress = 0
				item.State = qbittorrent.QbtDownloading
			} else {
				log.Debugf("%s: 等待: %v (%s)", name, int((item.Progress+0.005)*100), item.State)
			}
		case qbittorrent.QbtDownloading:
			item.Progress += 0.2
			if item.Progress >= 1.005 {
				item.Progress = 0
				item.State = qbittorrent.QbtUploading
			} else {
				log.Debugf("%s: 下载: %v (%s)", name, int((item.Progress+0.005)*100), item.State)
			}
		case qbittorrent.QbtUploading:
			item.Progress += 0.2
			if item.Progress >= 1.005 {
				item.State = qbittorrent.QbtCheckingUP
			} else {
				log.Debugf("%s: 做种: %v (%s)", name, int((item.Progress+0.005)*100), item.State)
			}
		default:
			log.Debugf("%s: %s", name, item.State)
		}
	}
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	_ = utils.CreateMutiDir(DownloadPath)
	_ = utils.CreateMutiDir(SavePath)
	out = bytes.NewBuffer(nil)

	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	wg := sync.WaitGroup{}
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	dirdb.Init(&dirdb.Options{
		DefaultExt: []string{".a_json", ".s_json", ".e_json"}, // anime, season
	})
	renamePlugin = renamer.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "rename/builtin_rename.py",
	})

	rename = wire.GetRenamer(
		&renamer.Options{
			WG:            &wg,
			RefreshSecond: 1,
		},
		&models.Plugin{
			Enable: true,
			Type:   "python",
			File:   "rename/builtin_rename.py",
		},
	)
	db = cache.NewBolt()
	db.Open("data/test.db")

	var err error
	callback := &clientnotifier.Callback{}
	dbInst, err := database.NewDatabase(&database.Options{
		SavePath: SavePath,
	}, db)
	if err != nil {
		panic(err)
	}
	dbs = clientnotifier.NewNotifier(&clientnotifier.Options{
		DownloadPath: DownloadPath,
		SavePath:     SavePath,
		Rename:       "link_delete",
		Callback:     callback,
	}, dbInst, rename)

	qbt = &qbittorrent.ClientMock{}
	qbt.MockInit(qbittorrent.ClientMockOptions{
		DownloadPath: DownloadPath,
		UpdateList:   defaultUpdateList,
		Ctx:          context.Background(),
	})

	mgr = downloader.NewManager(&downloader.Options{
		RefreshSecond:          1,
		Category:               "AnimeGoTest",
		WG:                     &wg,
		AllowDuplicateDownload: false,
		Tag:                    "",
	}, qbt, dbs)
	callback.Renamed = func(data any) error {
		return mgr.Delete(data.(string))
	}

	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func download(name string, season int, ep []int) (files []string, fullname, hash string) {
	hash = utils.MD5([]byte(fmt.Sprintf("%v%v%v", name, season, ep)))
	anime := &models.AnimeEntity{
		NameCN: name,
		Season: season,
		Torrent: &models.AnimeTorrent{
			Hash: hash,
		},
	}
	anime.Ep = make([]*models.AnimeEpEntity, 0, len(ep))
	for _, e := range ep {
		anime.Ep = append(anime.Ep, &models.AnimeEpEntity{
			Type: models.AnimeEpNormal,
			Ep:   e,
			Src:  fmt.Sprintf("%s/src_%d.mp4", name, e),
		})
	}
	fullname = anime.FullName()
	qbt.MockAddName(fullname, hash, anime.FilePathSrc())
	err := mgr.Download(anime)

	if err != nil {
		if !exceptions.IsExist(err) {
			panic(err)
		}
	}
	srcFiles := anime.FilePathSrc()
	files = make([]string, len(anime.Ep))
	log.ReInt(&log.Options{
		Debug: false,
	})
	for i := range anime.Ep {
		res, _ := renamePlugin.Rename(anime, i, srcFiles[i])
		files[i] = path.Join(SavePath, res.Filename)
	}
	log.ReInt(&log.Options{
		Debug: true,
	})
	return
}

func initTest(clean bool) (*sync.WaitGroup, func()) {
	if clean {
		_ = os.RemoveAll("data")
	}
	wg := sync.WaitGroup{}
	mgr.WG = &wg
	rename.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())

	rename.Init()
	rename.Start(ctx)
	dbs.Init()
	_ = dbs.Scan()

	qbt.Start()
	mgr.Init()
	mgr.Start(ctx)
	return &wg, cancel
}

func TestManager_Start(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)

	go func() {
		time.Sleep(14 * time.Second)
		download("test", 2, []int{1, 2, 3})
		time.Sleep(2 * time.Second)
		cancel()
	}()
	file1, _, _ := download("test", 2, []int{1, 2, 3})
	file2, _, _ := download("test", 2, []int{1, 2, 4})
	wg.Wait()
	time.Sleep(1 * time.Second)
	exist := dbs.IsExist(&models.AnimeEntity{
		NameCN: "test",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
			{Type: models.AnimeEpNormal, Ep: 4},
		},
	})
	assert.Equal(t, exist, true)
	for _, f := range file1 {
		assert.FileExists(t, f)
	}
	for _, f := range file2 {
		assert.FileExists(t, f)
	}
}

func TestManager_ReStartOnDownload(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)
	file1, _, _ := download("test", 2, []int{1, 2, 3})
	go func() {
		time.Sleep(3*time.Second + 500*time.Millisecond)
		cancel()
	}()
	wg.Wait()
	time.Sleep(1*time.Second + 500*time.Millisecond)

	wg, cancel = initTest(false)
	download("test", 2, []int{1, 2, 3})
	go func() {
		time.Sleep(11 * time.Second)
		cancel()
	}()
	wg.Wait()
	exist := dbs.IsExist(&models.AnimeEntity{
		NameCN: "test",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
		},
	})
	assert.Equal(t, exist, true)
	time.Sleep(1 * time.Second)
	for _, f := range file1 {
		assert.FileExists(t, f)
	}
}

func TestManager_ReStartOnSeed(t *testing.T) {
	out.Reset()
	wg, cancel := initTest(true)
	file1, _, _ := download("test", 2, []int{1, 2, 3})
	go func() {
		time.Sleep(8*time.Second + 500*time.Millisecond)
		cancel()
	}()
	wg.Wait()
	time.Sleep(1*time.Second + 500*time.Millisecond)

	wg, cancel = initTest(false)
	download("test", 2, []int{1, 2, 3})
	go func() {
		time.Sleep(6 * time.Second)
		cancel()
	}()
	wg.Wait()
	time.Sleep(1 * time.Second)
	exist := dbs.IsExist(&models.AnimeEntity{
		NameCN: "test",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
		},
	})
	assert.Equal(t, exist, true)
	for _, f := range file1 {
		assert.FileExists(t, f)
	}
}

func TestManager_StartWait(t *testing.T) {
	out.Reset()
	qbt.MockInit(qbittorrent.ClientMockOptions{
		DownloadPath: DownloadPath,
		UpdateList:   waitUpdateList,
		Ctx:          context.Background(),
	})
	wg, cancel := initTest(true)

	go func() {
		time.Sleep(14 * time.Second)
		download("test", 2, []int{1, 2, 3})
		time.Sleep(2 * time.Second)
		cancel()
	}()
	file1, _, _ := download("test", 2, []int{1, 2, 3})
	file2, _, _ := download("test", 2, []int{1, 2, 4})
	wg.Wait()
	time.Sleep(1 * time.Second)
	exist := dbs.IsExist(&models.AnimeEntity{
		NameCN: "test",
		Season: 2,
		Ep: []*models.AnimeEpEntity{
			{Type: models.AnimeEpNormal, Ep: 1},
			{Type: models.AnimeEpNormal, Ep: 2},
			{Type: models.AnimeEpNormal, Ep: 3},
			{Type: models.AnimeEpNormal, Ep: 4},
		},
	})
	assert.Equal(t, exist, true)
	for _, f := range file1 {
		assert.FileExists(t, f)
	}
	for _, f := range file2 {
		assert.FileExists(t, f)
	}
}

// TODO: 补充单测: 各种状态下重启

package database_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/utils"

	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
)

var (
	rename       *renamer.Manager
	renamePlugin *renamerPlugin.Rename
	db           *cache.Bolt
	out          *bytes.Buffer

	dbManager *database.Database
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
	wg := sync.WaitGroup{}
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	dirdb.Init(&dirdb.Options{
		DefaultExt: []string{".a_json", ".s_json"}, // anime, season
	})
	database.Init(&database.Options{
		DownloaderConf: database.DownloaderConf{
			DownloadPath: DownloadPath,
			SavePath:     SavePath,
			Rename:       "wait_move",
		},
	})
	renamer.Init(&renamer.Options{
		WG:                &wg,
		UpdateDelaySecond: 1,
	})
	renamePlugin = renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "rename/builtin_rename.py",
	})
	rename = renamer.NewManager(renamePlugin)
	db = cache.NewBolt()
	db.Open("data/test.db")
	var err error
	dbManager, err = database.NewDatabase(db, rename)
	if err != nil {
		panic(err)
	}

	m.Run()
	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}
func initTest() (*sync.WaitGroup, func()) {
	wg := sync.WaitGroup{}
	renamer.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())
	rename.Start(ctx)
	return &wg, cancel
}

func AddItem(name string, season int, ep []int) (hash string) {
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
	err := dbManager.Add(anime)
	if err != nil {
		panic(err)
	}
	for _, e := range anime.Ep {
		f := path.Join(DownloadPath, e.Src)
		err := utils.CreateMutiDir(path.Dir(f))
		if err != nil {
			panic(err)
		}
		err = os.WriteFile(f, []byte{}, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return hash
}

func TestOnDownloadStart(t *testing.T) {
	database.CacheMode = false
	wg, cancel := initTest()
	hash := AddItem("动画1", 2, []int{1, 2, 3})
	dbManager.OnDownloadStart([]models.ClientEvent{
		{Hash: hash},
	})
	dbManager.OnDownloadSeeding([]models.ClientEvent{
		{Hash: hash},
	})
	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	time.Sleep(1 * time.Second)
	dbManager.OnDownloadComplete([]models.ClientEvent{
		{Hash: hash},
	})
	wg.Wait()
}

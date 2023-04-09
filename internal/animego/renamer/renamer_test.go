package renamer_test

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/renamer"
	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"
)

var (
	ctx, cancel = context.WithCancel(context.Background())
	wg          = sync.WaitGroup{}
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.MkdirAll(DownloadPath, os.ModePerm)
	_ = os.MkdirAll(SavePath, os.ModePerm)
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	plugin.Init(&plugin.Options{
		Path:  "plugin/testdata",
		Debug: true,
	})
	m.Run()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func rename(r *renamer.Manager, state <-chan models.TorrentState, mode string, anime *models.AnimeEntity) string {
	srcs := anime.FilePathSrc()
	for _, s := range srcs {
		_ = os.WriteFile(path.Join(DownloadPath, s), []byte{}, os.ModePerm)
	}
	r.AddRenameTask(&models.RenameOptions{
		Entity: anime,
		SrcDir: DownloadPath,
		DstDir: SavePath,
		Mode:   mode,
		State:  state,
		RenameCallback: func(result *models.RenameResult) {
			d, _ := json.Marshal(result)
			fmt.Println(string(d))
		},
		CompleteCallback: func(result *models.RenameResult) {
			fmt.Println("exit", anime.DirName())
		},
	})
	return xpath.Join(SavePath, anime.DirName(), anime.FileName(0)+".mp4")
}

func Rename1(r *renamer.Manager, state <-chan models.TorrentState) string {
	mode := "link_delete"
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ドラえもん",
		NameCN:  "动画1",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep: []*models.AnimeEpEntity{
			{Ep: 712, Src: "src_712.mp4"},
			{Ep: 713, Src: "src_713.mp4"},
			{Ep: 714, Src: "src_714.mp4"},
		},
		MikanID: 681,
	}
	return rename(r, state, mode, anime)
}

func Rename2(r *renamer.Manager, state <-chan models.TorrentState) string {
	mode := "wait_move"
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "动画2",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep: []*models.AnimeEpEntity{
			{Ep: 1026, Src: "src_1026.mp4"},
		},
		MikanID: 228,
	}
	return rename(r, state, mode, anime)
}

func Rename3(r *renamer.Manager, state <-chan models.TorrentState) string {
	mode := "move"
	anime := &models.AnimeEntity{
		ID:      18692,
		Name:    "ONE PIECE",
		NameCN:  "动画3",
		AirDate: "2005-04-15",
		Eps:     0,
		Season:  1,
		Ep: []*models.AnimeEpEntity{
			{Ep: 996, Src: "src_996.mp4"},
		},
		MikanID: 228,
	}
	return rename(r, state, mode, anime)
}

func TestRenamer_Start(t *testing.T) {
	renamer.Init(&renamer.Options{
		WG:                &wg,
		UpdateDelaySecond: 3,
	})
	p := renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "rename.py",
	})
	r := renamer.NewManager(p)
	state1 := make(chan models.TorrentState)
	f1 := Rename1(r, state1)
	state2 := make(chan models.TorrentState)
	f2 := Rename2(r, state2)
	state3 := make(chan models.TorrentState)
	f3 := Rename3(r, state3)

	r.Start(ctx)
	time.Sleep(500 * time.Millisecond)
	go func() {
		state1 <- downloader.StateSeeding
		state1 <- downloader.StateComplete
	}()
	go func() {
		state2 <- downloader.StateSeeding
		state2 <- downloader.StateComplete
	}()
	go func() {
		state3 <- downloader.StateSeeding
		state3 <- downloader.StateComplete
	}()

	go func() {
		time.Sleep(10 * time.Second)
		cancel()
	}()
	wg.Wait()

	assert.FileExists(t, f1)
	assert.FileExists(t, f2)
	assert.FileExists(t, f3)
}

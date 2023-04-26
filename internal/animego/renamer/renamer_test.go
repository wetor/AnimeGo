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
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func rename(r *renamer.Manager, state []chan models.TorrentState, mode string, anime *models.AnimeEntity) []string {
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
			fmt.Println("complete", anime.DirName())
		},
	})
	dst := anime.FilePath()
	result := make([]string, len(dst))
	for i := range result {
		result[i] = xpath.Join(SavePath, dst[i])
	}
	return result
}

func makeChans(size int) []chan models.TorrentState {
	state := make([]chan models.TorrentState, size)
	for i := range state {
		state[i] = make(chan models.TorrentState)
	}
	return state
}

func Rename1(r *renamer.Manager) ([]string, []chan models.TorrentState) {
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
	state := makeChans(len(anime.Ep))
	return rename(r, state, mode, anime), state
}

func Rename2(r *renamer.Manager) ([]string, []chan models.TorrentState) {
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
	state := makeChans(len(anime.Ep))
	return rename(r, state, mode, anime), state
}

func Rename3(r *renamer.Manager) ([]string, []chan models.TorrentState) {
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
	state := makeChans(len(anime.Ep))
	return rename(r, state, mode, anime), state
}

func TestRenamer_Start(t *testing.T) {
	renamer.Init(&renamer.Options{
		WG:                &wg,
		UpdateDelaySecond: 1,
	})
	p := renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_rename.py",
	})
	r := renamer.NewManager(p)
	f1, state1 := Rename1(r)
	f2, state2 := Rename2(r)
	f3, state3 := Rename3(r)

	r.Start(ctx)
	time.Sleep(500 * time.Millisecond)
	go func() {
		for i := range state1 {
			state1[i] <- downloader.StateSeeding
			state1[i] <- downloader.StateComplete
		}
	}()
	go func() {
		for i := range state2 {
			state2[i] <- downloader.StateSeeding
			state2[i] <- downloader.StateComplete
		}
	}()
	go func() {
		for i := range state3 {
			state3[i] <- downloader.StateSeeding
			state3[i] <- downloader.StateComplete
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)
		cancel()
	}()
	wg.Wait()

	for i := range f1 {
		assert.FileExists(t, f1[i])
	}
	for i := range f2 {
		assert.FileExists(t, f2[i])
	}
	for i := range f3 {
		assert.FileExists(t, f3[i])
	}
}

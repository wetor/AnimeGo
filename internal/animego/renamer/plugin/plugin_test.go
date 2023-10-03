package plugin_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	renamerPlugin "github.com/wetor/AnimeGo/internal/animego/renamer/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	plugin.Init(&plugin.Options{
		Path:  "testdata",
		Debug: true,
	})
	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestRename_Rename(t *testing.T) {
	type args struct {
		anime *models.AnimeEntity
		index int
	}
	tests := []struct {
		name string
		args args
		want *models.RenameResult
	}{
		// TODO: Add test cases.
		{
			args: args{
				anime: &models.AnimeEntity{
					ID:     18692,
					Name:   "ドラえもん",
					NameCN: "动画1",
					Season: 1,
					Ep: []*models.AnimeEpEntity{
						{Ep: 712, Src: "src_712.mp4", Type: models.AnimeEpNormal},
					},
				},
				index: 0,
			},
			want: &models.RenameResult{Index: 0, Filename: "动画1/S01/E712.mp4", AnimeDir: "动画1", SeasonDir: "动画1/S01", Scrape: true},
		},
		{
			args: args{
				anime: &models.AnimeEntity{
					ID:     18692,
					Name:   "ドラえもん",
					NameCN: "动画1",
					Season: 1,
					Ep: []*models.AnimeEpEntity{
						{Ep: 712, Src: "src_712.mp4", Type: models.AnimeEpNormal},
					},
				},
				index: 0,
			},
			want: &models.RenameResult{Index: 0, Filename: "动画1/S01/E712.mp4", AnimeDir: "动画1", SeasonDir: "动画1/S01", Scrape: true},
		},
		{
			args: args{
				anime: &models.AnimeEntity{
					ID:     18692,
					Season: 2,
					Ep: []*models.AnimeEpEntity{
						{Ep: 0, Src: "src_1.mp4", Type: models.AnimeEpUnknown},
					},
				},
				index: 0,
			},
			want: &models.RenameResult{Index: 0, Filename: "18692/S02/src_1.mp4", AnimeDir: "18692", SeasonDir: "18692/S02", Scrape: true},
		},
	}
	p := renamerPlugin.NewRenamePlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_rename.py",
	})

	for _, tt := range tests {
		filename := tt.args.anime.FilePathSrc()
		t.Run(tt.name, func(t *testing.T) {
			tt.args.anime.Default()
			result, _ := p.Rename(tt.args.anime, tt.args.index, filename[tt.args.index])
			assert.Equalf(t, tt.want, result, "Rename(%v, %v)", tt.args.anime, filename)
		})
	}
}

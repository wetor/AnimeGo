package plugin_test

import (
	"fmt"
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
	fmt.Println("end")
}

func TestRename_Rename(t *testing.T) {
	type args struct {
		anime    *models.AnimeEntity
		filename string
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
					Ep:     712,
				},
				filename: "test/path/v.mp4",
			},
			want: &models.RenameResult{Filepath: "动画1/S01/E712.mp4", TVShowDir: "动画1"},
		},
		{
			args: args{
				anime: &models.AnimeEntity{
					ID:     18692,
					Name:   "ドラえもん",
					NameCN: "动画1",
					Season: 1,
					Ep:     712,
				},
				filename: "test/path/v.mp4",
			},
			want: &models.RenameResult{Filepath: "动画1/S01/E712.mp4", TVShowDir: "动画1"},
		},
		{
			args: args{
				anime: &models.AnimeEntity{
					ID:     18692,
					Season: 2,
					Ep:     1,
				},
				filename: "test/path/v.mp4",
			},
			want: &models.RenameResult{Filepath: "18692/S02/E1.mp4", TVShowDir: "18692"},
		},
	}
	p := renamerPlugin.NewRenamePlugin([]models.Plugin{
		{
			Enable: true,
			Type:   "python",
			File:   "rename.py",
		},
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, p.Rename(tt.args.anime, tt.args.filename), "Rename(%v, %v)", tt.args.anime, tt.args.filename)
		})
	}
}

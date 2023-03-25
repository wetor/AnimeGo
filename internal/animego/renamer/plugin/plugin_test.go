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
		Path:  "../../../../assets/plugin",
		Debug: true,
	})
	m.Run()
	fmt.Println("end")
}

func TestRename_Rename(t *testing.T) {
	p := renamerPlugin.NewRenamePlugin([]models.Plugin{
		{
			Enable: true,
			Type:   "python",
			File:   "rename/builtin_rename.py",
		},
	})
	dst := p.Rename(&models.AnimeEntity{
		ID:     18692,
		Name:   "ドラえもん",
		NameCN: "动画1",
		Season: 1,
		Ep:     712,
	}, "test/path/v.mp4")
	assert.Equal(t, dst, "动画1/S01/E712.mp4")

	dst = p.Rename(&models.AnimeEntity{
		ID:     18692,
		Name:   "ドラえもん",
		Season: 1,
		Ep:     7121,
	}, "test/path/v.mp4")
	assert.Equal(t, dst, "ドラえもん/S01/E7121.mp4")

	dst = p.Rename(&models.AnimeEntity{
		ID:     18692,
		Season: 2,
		Ep:     1,
	}, "test/path/v.mp4")
	assert.Equal(t, dst, "18692/S02/E1.mp4")
}

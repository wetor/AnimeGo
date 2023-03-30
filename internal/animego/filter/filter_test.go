package filter_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
)

var (
	mgr *filter.Manager
	ctx = context.Background()
)

type MockManager struct {
}

func (m *MockManager) Download(anime any) {
	d, _ := json.Marshal(anime)
	fmt.Println(string(d))
}

type MockFeed struct {
}

func (m *MockFeed) Parse(opt *models.AnimeParseOptions) *models.AnimeEntity {
	return &models.AnimeEntity{
		ID:           975,
		ThemoviedbID: 37854,
		MikanID:      228,
		Name:         "ONE PIECE",
		NameCN:       "海贼王",
		Season:       opt.Season,
		Ep:           opt.Ep,
		Eps:          1079,
		AirDate:      "1999-10-20",
	}
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	plugin.Init(&plugin.Options{
		Path:  "../../../assets/plugin",
		Debug: true,
	})

	filter.Init(&filter.Options{
		DelaySecond: 1,
	})

	mgr = filter.NewManager(&MockFeed{}, &MockManager{})
	m.Run()

	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestManager_Update(t *testing.T) {
	items := []*models.FeedItem{
		{
			Url:      "url1",
			Name:     "OPFans枫雪动漫][ONE PIECE 海贼王][第1029话][1080p][周日版][MP4][简体]",
			Download: "download1",
		},
	}
	mgr.Update(ctx, items)
}

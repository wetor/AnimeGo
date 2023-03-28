package filter_test

import (
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/pkg/json"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
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
	return &models.AnimeEntity{}
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
			Name:     "name1",
			Download: "download1",
		},
	}
	mgr.Update(ctx, items)
}

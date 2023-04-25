package filter_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"testing"

	"github.com/brahma-adshonor/gohook"

	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
)

var (
	mgr *filter.Manager
	ctx = context.Background()
)

func HookGetWriter(uri string, w io.Writer) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := path.Base(uri)
	jsonData, err := os.ReadFile(path.Join("testdata", id))
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

type MockManager struct {
}

func (m *MockManager) Download(anime any) {
	d, _ := json.Marshal(anime)
	fmt.Println(string(d))
}

type MockParser struct {
}

func (m *MockParser) Parse(opt *models.ParseOptions) *models.AnimeEntity {
	return &models.AnimeEntity{
		ID:           975,
		ThemoviedbID: 37854,
		MikanID:      228,
		Season:       1,
		Name:         "ONE PIECE",
		NameCN:       "海贼王",
		Eps:          1079,
		AirDate:      "1999-10-20",
		Ep: []*models.AnimeEpEntity{
			{
				Ep:  109,
				Src: "src_109.mp4",
			},
			{
				Ep:  110,
				Src: "src_110.mp4",
			},
		},
	}
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	_ = gohook.Hook(request.GetWriter, HookGetWriter, nil)

	plugin.Init(&plugin.Options{
		Path:  "../../../assets/plugin",
		Debug: true,
	})

	filter.Init(&filter.Options{
		DelaySecond: 1,
	})

	mgr = filter.NewManager(&MockManager{}, &MockParser{})
	m.Run()

	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestManager_Update(t *testing.T) {
	items := []*models.FeedItem{
		{
			Url:      "url1",
			Name:     "[轻之国度字幕组][想要成为影之实力者！/我想成为影之强者！][12][GB][720P][MP4][附小剧场]",
			Download: "https://mikanani.me/Download/20221223/8b0b9621f7f7b8425ed0f6162d03b92c93db2270.torrent",
		},
		{
			Url:      "url2",
			Name:     "[轻之国度字幕组][想要成为影之实力者！/我想成为影之强者！][12][GB][720P][MP4][附小剧场]",
			Download: "magnet:?xt=urn:btih:4c81fc90f8db37eae70a29a82e7abf8d8f1867c2&dn=example+file&tr=udp%3A%2F%2Ftracker.example.com%3A80",
		},
		{
			Url:      "https://mikanani.me/Home/Episode/1069f01462c90b1065f4fc5576529422451a90a9",
			Name:     "[猎户不鸽压制] 万事屋斋藤先生转生异世界 / 斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku [01-12] [合集] [WebRip 1080p] [简中内嵌] [H265 AAC] [2023年1月番]",
			Download: "https://mikanani.me/Download/20230328/1069f01462c90b1065f4fc5576529422451a90a9.torrent",
		},
	}
	mgr.Update(ctx, items)
}

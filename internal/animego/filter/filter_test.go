package filter_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "filter"

var (
	mgr *filter.Manager
	ctx = context.Background()
)

type DownloaderMock struct {
}

func (m *DownloaderMock) Download(anime *models.AnimeEntity) error {
	d, _ := json.Marshal(anime)
	fmt.Println(string(d))
	return nil
}

type ParserMock struct {
}

func (m *ParserMock) Parse(opt *models.ParseOptions) (*models.AnimeEntity, error) {
	switch opt.Title {
	case "err_parse":
		return nil, &exceptions.ErrParseFailed{}
	}
	return &models.AnimeEntity{
		ID:           975,
		ThemoviedbID: 37854,
		MikanID:      228,
		Season:       1,
		Name:         opt.Title,
		NameCN:       opt.Title,
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
	}, nil
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	plugin.Init(&plugin.Options{
		Path:  test.GetDataPath(testdata, ""),
		Debug: true,
	})

	filter.Init(&filter.Options{
		DelaySecond: 1,
	})

	mgr = filter.NewManager(&DownloaderMock{}, &ParserMock{})
	mgr.Add(&models.Plugin{
		Enable: true,
		Type:   "py",
		File:   "default.py",
	})
	m.Run()

	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestManager_Update(t *testing.T) {
	type args struct {
		ctx        context.Context
		items      []*models.FeedItem
		skipFilter bool
		skipDelay  bool
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				ctx: ctx,
				items: []*models.FeedItem{
					{
						MikanUrl:   "url1",
						Name:       "[轻之国度字幕组][想要成为影之实力者！/我想成为影之强者！][12][GB][720P][MP4][附小剧场]",
						TorrentUrl: "https://mikanani.me/Download/20221223/8b0b9621f7f7b8425ed0f6162d03b92c93db2270.torrent",
					},
					{
						MikanUrl:   "url2",
						Name:       "[轻之国度字幕组][想要成为影之实力者！/我想成为影之强者！][12][GB][720P][MP4][附小剧场]",
						TorrentUrl: "magnet:?xt=urn:btih:4c81fc90f8db37eae70a29a82e7abf8d8f1867c2&dn=example+file&tr=udp%3A%2F%2Ftracker.example.com%3A80",
					},
					{
						MikanUrl:   "https://mikanani.me/Home/Episode/1069f01462c90b1065f4fc5576529422451a90a9",
						Name:       "[猎户不鸽压制] 万事屋斋藤先生转生异世界 / 斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku [01-12] [合集] [WebRip 1080p] [简中内嵌] [H265 AAC] [2023年1月番]",
						TorrentUrl: "https://mikanani.me/Download/20230328/1069f01462c90b1065f4fc5576529422451a90a9.torrent",
					},
				},
				skipFilter: false,
				skipDelay:  true,
			},
		},
		{
			name: "no_item",
			args: args{
				ctx:        ctx,
				items:      []*models.FeedItem{},
				skipFilter: false,
				skipDelay:  true,
			},
		},
		{
			name: "err_parse",
			args: args{
				ctx: ctx,
				items: []*models.FeedItem{
					{
						MikanUrl:   "err_parse",
						Name:       "err_parse",
						TorrentUrl: "err_parse",
					},
				},
				skipFilter: true,
				skipDelay:  true,
			},
		},
		{
			name: "delay",
			args: args{
				ctx: ctx,
				items: []*models.FeedItem{
					{
						MikanUrl:   "delay",
						Name:       "delay",
						TorrentUrl: "delay",
					},
				},
				skipFilter: true,
				skipDelay:  false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mgr.Update(tt.args.ctx, tt.args.items, tt.args.skipFilter, tt.args.skipDelay)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

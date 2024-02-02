package anisource_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/anisource/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource/themoviedb"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/pkg/request"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/wire"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

var (
	mikanSource api.AniSource
	ctx, cancel = context.WithCancel(context.Background())
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	b := cache.NewBolt()
	b.Open("data/bolt.db")

	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open(test.GetDataPath("", "bolt_sub.bolt"))
	bangumiCache.Add("bangumi_sub")
	mutex := sync.Mutex{}

	mikanSource = wire.GetMikan(&mikan.Options{
		Cache:     b,
		CacheTime: int64(7 * 24 * 60 * 60),
	}, &bangumi.Options{
		Cache:            b,
		CacheTime:        int64(7 * 24 * 60 * 60),
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &mutex,
	}, &themoviedb.Options{
		Cache:     b,
		CacheTime: int64(7 * 24 * 60 * 60),
	})

	bangumiHost := test.MockBangumiStart(ctx)
	mikanHost := test.MockMikanStart(ctx)
	themoviedbHost := test.MockThemoviedbStart(ctx)
	request.Init(&request.Options{
		Host: map[string]*request.HostOptions{
			constant.BangumiHost: {
				Redirect: bangumiHost,
			},
			constant.MikanHost: {
				Redirect: mikanHost,
				Cookie: map[string]string{
					constant.MikanAuthCookie: "MikanAuthCookie",
				},
			},
			constant.ThemoviedbHost: {
				Redirect: themoviedbHost,
				Params: map[string]string{
					constant.ThemoviedbApiKey: "123456",
				},
			},
		},
	})

	m.Run()

	cancel()
	b.Close()
	bangumiCache.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestMikan_Parse(t *testing.T) {
	type args struct {
		opts *models.AnimeParseOptions
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantAnime  *models.AnimeEntity
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name: "海贼王",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "https://mikanani.me/Home/Episode/18b60d48a72c603b421468aade7fdd0868ff2f2f",
				},
				name: "OPFans枫雪动漫][ONE PIECE 海贼王][第1029话][1080p][周日版][MP4][简体] [299.5MB]",
			},
			wantAnime: &models.AnimeEntity{ID: 975, ThemoviedbID: 37854, MikanID: 228, Name: "ONE PIECE", NameCN: "海贼王", Season: 1, Eps: 1079, AirDate: "1999-10-20"},
		},
		{
			name: "欢迎来到实力至上主义的教室 第二季",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "https://mikanani.me/Home/Episode/8849c25e05d6e2623b5333bc78d3a489a9b1cc59",
				},
				name: "[ANi] Classroom of the Elite S2 - 欢迎来到实力至上主义的教室 第二季 - 07 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4] [254.26 MB]",
			},
			wantAnime: &models.AnimeEntity{ID: 371546, ThemoviedbID: 72517, MikanID: 2775, Name: "ようこそ実力至上主義の教室へ 2nd Season", NameCN: "欢迎来到实力至上主义教室 第二季", Season: 2, Eps: 13, AirDate: "2022-07-04"},
		},
		{
			name: "想要成为影之实力者",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "https://mikanani.me/Home/Episode/dcc28079dfda415cdcdf46159aad0fa94f1a2f11",
				},
				name: "[LoliHouse] 想要成为影之实力者 / 我想成为影之强者 / Kage no Jitsuryokusha ni Naritakute! - 19 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕]",
			},
			wantAnime: &models.AnimeEntity{ID: 329114, ThemoviedbID: 119495, MikanID: 2822, Name: "陰の実力者になりたくて！", NameCN: "想要成为影之实力者！", Season: 1, Eps: 20, AirDate: "2022-10-05"},
		},
		{
			name: "AnimeParseOverride",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "https://mikanani.me/Home/Episode/dcc28079dfda415cdcdf46159aad0fa94f1a2f11",
					AnimeParseOverride: &models.AnimeParseOverride{
						MikanID:      114,
						BangumiID:    514,
						ThemoviedbID: 1919,
						Name:         "AnimeParseOverride",
						NameCN:       "AnimeParseOverrideCN",
						AirDate:      "2022-10-05",
						Season:       1,
						Eps:          20,
					},
				},
				name: "AnimeParseOverride",
			},
			wantAnime: &models.AnimeEntity{ID: 514, ThemoviedbID: 1919, MikanID: 114, Name: "AnimeParseOverride", NameCN: "AnimeParseOverrideCN", Season: 1, Eps: 20, AirDate: "2022-10-05"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAnime, err := mikanSource.Parse(tt.args.opts)
			assert.NoError(t, err)
			assert.Equalf(t, tt.wantAnime, gotAnime, "Parse(%v)", tt.args.opts)
		})
	}
}

func TestMikan_Parse_Failed(t *testing.T) {
	type args struct {
		opts *models.AnimeParseOptions
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantAnime  *models.AnimeEntity
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name: "err_mikan",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "err_mikan",
				},
				name: "err_mikan",
			},
			wantAnime:  nil,
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "请求 err_mikan 失败",
		},
		{
			name: "err_bangumi",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "err_bangumi",
				},
				name: "err_bangumi",
			},
			wantAnime:  nil,
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "请求 err_bangumi 失败",
		},
		{
			name: "err_themoviedb_search",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "err_themoviedb_search",
				},
				name: "err_themoviedb_search",
			},
			wantAnime:  nil,
			wantErr:    &exceptions.ErrRequest{},
			wantErrStr: "请求 err_themoviedb_search 失败",
		},
		{
			name: "err_themoviedb_get",
			args: args{
				opts: &models.AnimeParseOptions{
					Input: "err_themoviedb_get",
				},
				name: "err_themoviedb_get",
			},
			wantAnime: &models.AnimeEntity{ID: 1919, ThemoviedbID: 666, MikanID: 1919, Name: "err_themoviedb_get", NameCN: "err_themoviedb_get", Season: 0, Ep: nil, Eps: 2, AirDate: "1919-05-14"},
		},
	}
	save := mikanSource
	mikanSource = MikanSourceMock()
	defer func() {
		mikanSource = save
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := mikanSource.Parse(tt.args.opts)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantAnime, gotEntity, "Parse(%v)", tt.args.opts)
			}
		})
	}
}

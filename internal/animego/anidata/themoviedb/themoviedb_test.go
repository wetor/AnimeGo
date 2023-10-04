package themoviedb_test

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/wetor/AnimeGo/pkg/xpath"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "themoviedb"

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	db := cache.NewBolt()
	db.Open("data/bolt.db")
	anidata.Init(&anidata.Options{Cache: db})
	request.Init(&request.Options{
		Debug: true,
	})
	test.HookGet(testdata, func(uri string) string {
		u, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		id := u.Query().Get("with_text_query")
		if len(id) == 0 {
			id = path.Base(xpath.P(u.Path))
		}
		return id
	})
	defer test.UnHook()
	m.Run()

	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestThemoviedb_Get_GetCache(t *testing.T) {
	type args struct {
		name    string
		airDate string
	}
	tests := []struct {
		name           string
		args           args
		wantID         int
		wantSeasonInfo *themoviedb.SeasonInfo
		wantErr1       error
		wantErr1Str    string
		wantErr2       error
		wantErr2Str    string
	}{
		// TODO: Add test cases.
		{
			name:           "海贼王",
			args:           args{name: "ONE PIECE", airDate: "1999-10-20"},
			wantID:         37854,
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 1, AirDate: "1999-10-20", EpID: 49188, EpName: "", Ep: 0, Eps: 61},
		},
		{
			name:           "在地下城寻求邂逅是否搞错了什么 Ⅳ 新章 迷宫篇",
			args:           args{name: "ダンジョンに出会いを求めるのは間違っているだろうか Ⅳ 新章 迷宮篇", airDate: "2022-07-21"},
			wantID:         62745,
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 4, AirDate: "2022-07-23", EpID: 193725, EpName: "", Ep: 0, Eps: 22},
		},
		{
			name:           "来自深渊 烈日的黄金乡",
			args:           args{name: "メイドインアビス 烈日の黄金郷", airDate: "2022-07-06"},
			wantID:         72636,
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 2, AirDate: "2022-07-06", EpID: 204984, EpName: "", Ep: 0, Eps: 12},
		},
		{
			name:           "OVERLORD IV",
			args:           args{name: "オーバーロードIV", airDate: "2022-07-05"},
			wantID:         64196,
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 4, AirDate: "2022-07-05", EpID: 194087, EpName: "", Ep: 0, Eps: 13},
		},
		{
			name:           "福星小子",
			args:           args{name: "うる星やつら", airDate: "2022-10-14"},
			wantID:         154524,
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 1, AirDate: "2022-10-14", EpID: 237892, EpName: "", Ep: 0, Eps: 46},
		},
		{
			name:           "Mairimashita! Iruma-kun 3rd Season",
			args:           args{name: "Mairimashita! Iruma-kun 3rd Season", airDate: "2022-11-14"},
			wantID:         91801,
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 3, AirDate: "2022-10-08", EpID: 306624, EpName: "", Ep: 0, Eps: 21},
		},
		{
			name:        "err_search_not_found",
			args:        args{name: "番剧IV 2期 2nd Season 第二季", airDate: "1999-10-20"},
			wantErr1:    &exceptions.ErrThemoviedbSearchName{},
			wantErr1Str: "查询ThemoviedbID失败: 搜索番剧名失败",
		},
		{
			name:        "err_search_not_similar",
			args:        args{name: "err_search_not_similar", airDate: "1999-10-20"},
			wantErr1:    &exceptions.ErrThemoviedbMatchSeason{},
			wantErr1Str: "查询ThemoviedbID失败: 匹配季度信息失败，番剧名未找到",
		},
		{
			name:        "err_get_season_request",
			args:        args{name: "err_get_season_request", airDate: "1999-10-20"},
			wantID:      114513,
			wantErr2:    &exceptions.ErrRequest{},
			wantErr2Str: "获取Themoviedb信息失败: 请求 Themoviedb 失败",
		},
		{
			name:        "err_get_season_null",
			args:        args{name: "err_get_season_null", airDate: "1999-10-20"},
			wantID:      114514,
			wantErr2:    &exceptions.ErrThemoviedbMatchSeason{},
			wantErr2Str: "获取Themoviedb信息失败: 匹配季度信息失败，此番剧可能未开播",
		},
		{
			name:        "err_get_season_no_match",
			args:        args{name: "err_get_season_no_match", airDate: "1999-10-20"},
			wantID:      11451419,
			wantErr2:    &exceptions.ErrThemoviedbMatchSeason{},
			wantErr2Str: "获取Themoviedb信息失败: 匹配季度信息失败，此番剧可能未开播",
		},
	}
	tmdb := &themoviedb.Themoviedb{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := tmdb.SearchCache(tt.args.name, nil)
			if tt.wantErr1 != nil {
				assert.IsType(t, tt.wantErr1, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErr1Str)
				return
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantID, id, "SearchCache(%v)", tt.args.name)
			}

			gotSeasonInfo, err := tmdb.GetCache(id, tt.args.airDate)
			if tt.wantErr2 != nil {
				assert.IsType(t, tt.wantErr2, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErr2Str)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantSeasonInfo, gotSeasonInfo, "GetCache(%v, %v)", tt.args.name, tt.args.airDate)
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			id, err := tmdb.Search(tt.args.name, nil)
			if tt.wantErr1 != nil {
				assert.IsType(t, tt.wantErr1, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErr1Str)
				return
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantID, id, "Search(%v)", tt.args.name)
			}

			gotSeasonInfo, err := tmdb.Get(id, tt.args.airDate)
			if tt.wantErr2 != nil {
				assert.IsType(t, tt.wantErr2, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErr2Str)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantSeasonInfo, gotSeasonInfo, "Get(%v, %v)", tt.args.name, tt.args.airDate)
			}
		})
	}
}

package themoviedb_test

import (
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/xpath"
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
			id = xpath.Base(u.Path)
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

func TestThemoviedb_ParseCache(t *testing.T) {
	type args struct {
		name    string
		airDate string
	}
	tests := []struct {
		name           string
		args           args
		wantEntity     *themoviedb.Entity
		wantSeasonInfo *themoviedb.SeasonInfo
	}{
		// TODO: Add test cases.
		{
			name:           "海贼王",
			args:           args{name: "ONE PIECE", airDate: "1999-10-20"},
			wantEntity:     &themoviedb.Entity{ID: 37854, NameCN: "海贼王", Name: "ワンピース", AirDate: "1999-10-20"},
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 1, AirDate: "1999-10-20", EpID: 49188, EpName: "", Ep: 0, Eps: 61},
		},
		{
			name:           "在地下城寻求邂逅是否搞错了什么 Ⅳ 新章 迷宫篇",
			args:           args{name: "ダンジョンに出会いを求めるのは間違っているだろうか Ⅳ 新章 迷宮篇", airDate: "2022-07-21"},
			wantEntity:     &themoviedb.Entity{ID: 62745, NameCN: "在地下城寻求邂逅是否搞错了什么", Name: "ダンジョンに出会いを求めるのは間違っているだろうか", AirDate: "2015-04-04"},
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 4, AirDate: "2022-07-23", EpID: 193725, EpName: "", Ep: 0, Eps: 22},
		},
		{
			name:           "来自深渊 烈日的黄金乡",
			args:           args{name: "メイドインアビス 烈日の黄金郷", airDate: "2022-07-06"},
			wantEntity:     &themoviedb.Entity{ID: 72636, NameCN: "来自深渊", Name: "メイドインアビス", AirDate: "2017-07-07"},
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 2, AirDate: "2022-07-06", EpID: 204984, EpName: "", Ep: 0, Eps: 12},
		},
		{
			name:           "OVERLORD IV",
			args:           args{name: "オーバーロードIV", airDate: "2022-07-05"},
			wantEntity:     &themoviedb.Entity{ID: 64196, NameCN: "不死者之王", Name: "オーバーロード", AirDate: "2015-07-07"},
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 4, AirDate: "2022-07-05", EpID: 194087, EpName: "", Ep: 0, Eps: 13},
		},
		{
			name:           "福星小子",
			args:           args{name: "うる星やつら", airDate: "2022-10-14"},
			wantEntity:     &themoviedb.Entity{ID: 154524, NameCN: "福星小子", Name: "うる星やつら", AirDate: "2022-10-14"},
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 1, AirDate: "2022-10-14", EpID: 237892, EpName: "", Ep: 0, Eps: 46},
		},
		{
			name:           "Mairimashita! Iruma-kun 3rd Season",
			args:           args{name: "Mairimashita! Iruma-kun 3rd Season", airDate: "2022-11-14"},
			wantEntity:     &themoviedb.Entity{ID: 91801, NameCN: "入间同学入魔了！", Name: "魔入りました！入間くん", AirDate: "2019-10-05"},
			wantSeasonInfo: &themoviedb.SeasonInfo{Season: 3, AirDate: "2022-10-08", EpID: 306624, EpName: "", Ep: 0, Eps: 21},
		},
	}
	tmdb := &themoviedb.Themoviedb{}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			id, err := tmdb.SearchCache(tt.args.name)
			assert.NoError(t, err)
			gotSeasonInfo, err := tmdb.GetCache(id, tt.args.airDate)
			assert.NoError(t, err)
			assert.Equalf(t1, tt.wantEntity.ID, id, "ParseCache(%v, %v)", tt.args.name, tt.args.airDate)
			assert.Equalf(t1, tt.wantSeasonInfo, gotSeasonInfo, "ParseCache(%v, %v)", tt.args.name, tt.args.airDate)
		})
	}
}

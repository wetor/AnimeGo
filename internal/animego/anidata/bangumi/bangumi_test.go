package bangumi_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "bangumi"

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	request.Init(&request.Options{
		Debug: true,
	})
	test.HookAll(testdata, nil)
	defer test.UnHook()
	mutex := sync.Mutex{}

	db := cache.NewBolt()
	db.Open("data/bolt.db")

	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open(test.GetDataPath("", "bolt_sub.bolt"))
	anidata.Init(&anidata.Options{
		Cache:            db,
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &mutex,
	})

	m.Run()

	db.Close()
	bangumiCache.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestBangumi_Parse(t *testing.T) {
	type args struct {
		bangumiID int
		ep        int
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *bangumi.Entity
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &bangumi.Entity{ID: 253047, NameCN: "联盟空军航空魔法音乐队 光辉魔女", Name: "連盟空軍航空魔法音楽隊 ルミナスウィッチーズ", Eps: 12, AirDate: "2022-07-03", Type: 0, Platform: 0},
			wantErr:    false,
		},
		{
			name:       "欢迎来到实力至上主义的教室 第二季",
			args:       args{bangumiID: 371546, ep: 7},
			wantEntity: &bangumi.Entity{ID: 371546, NameCN: "欢迎来到实力至上主义教室 第二季", Name: "ようこそ実力至上主義の教室へ 2nd Season", Eps: 13, AirDate: "2022-07-04", Type: 0, Platform: 0},
			wantErr:    false,
		},
	}
	b := &bangumi.Bangumi{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity := b.Get(tt.args.bangumiID, nil)
			assert.Equalf(t, tt.wantEntity, gotEntity, "Get(%v)", tt.args.bangumiID)
		})
	}
}

func TestBangumi_ParseCache(t *testing.T) {
	type args struct {
		bangumiID int
		ep        int
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *bangumi.Entity
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &bangumi.Entity{ID: 253047, NameCN: "联盟空军航空魔法音乐队 光辉魔女", Name: "連盟空軍航空魔法音楽隊 ルミナスウィッチーズ", Eps: 12, AirDate: "2022-07-03", Type: 0, Platform: 0},
			wantErr:    false,
		},
		{
			name:       "欢迎来到实力至上主义的教室 第二季",
			args:       args{bangumiID: 371546, ep: 5},
			wantEntity: &bangumi.Entity{ID: 371546, NameCN: "欢迎来到实力至上主义教室 第二季", Name: "ようこそ実力至上主義の教室へ 2nd Season", Eps: 13, AirDate: "2022-07-04", Type: 0, Platform: 0},
			wantErr:    false,
		},
		{
			name:       "CLANNAD",
			args:       args{bangumiID: 51, ep: 5},
			wantEntity: &bangumi.Entity{ID: 51, NameCN: "CLANNAD", Name: "CLANNAD -クラナド-", Eps: 22, AirDate: "2007-10-04", Type: 2, Platform: 1},
			wantErr:    false,
		},
	}

	b := &bangumi.Bangumi{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity := b.GetCache(tt.args.bangumiID, nil).(*bangumi.Entity)
			assert.Equalf(t, tt.wantEntity, gotEntity, "GetCache(%v)", tt.args.bangumiID)
		})
	}
}

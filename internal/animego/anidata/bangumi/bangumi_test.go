package bangumi

import (
	"fmt"
	"sync"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	utils.CreateMutiDir("data")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	request.Init(&request.Options{})
	mutex := sync.Mutex{}

	db := cache.NewBolt()
	db.Open("data/bolt.db")

	bangumiCache := cache.NewBolt()
	bangumiCache.Open("testdata/bolt_sub.bolt")
	anidata.Init(&anidata.Options{
		Cache:            db,
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &mutex,
	})

	m.Run()
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
		wantEntity *Entity
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &Entity{Eps: 12},
			wantErr:    false,
		},
		{
			name:       "欢迎来到实力至上主义的教室 第二季",
			args:       args{bangumiID: 371546, ep: 7},
			wantEntity: &Entity{Eps: 13},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bangumi{}
			gotEntity := b.Parse(tt.args.bangumiID)
			if gotEntity.Eps != tt.wantEntity.Eps {
				t.Errorf("Parse() gotEntity = %v, want %v", gotEntity, tt.wantEntity)
			}
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
		wantEntity *Entity
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &Entity{Eps: 12},
			wantErr:    false,
		},
		{
			name:       "欢迎来到实力至上主义的教室 第二季",
			args:       args{bangumiID: 371546, ep: 5},
			wantEntity: &Entity{Eps: 13},
			wantErr:    false,
		},
		{
			name:       "CLANNAD",
			args:       args{bangumiID: 51, ep: 5},
			wantEntity: &Entity{Eps: 22},
			wantErr:    false,
		},
	}

	b := &Bangumi{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotEntity := b.ParseCache(tt.args.bangumiID)
			if gotEntity.Eps != tt.wantEntity.Eps {
				t.Errorf("Parse() gotEntity = %v, want %v", gotEntity, tt.wantEntity)
			}
		})
	}
}

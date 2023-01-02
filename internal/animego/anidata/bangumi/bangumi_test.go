package bangumi

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/cache"
	"go.uber.org/zap"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
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
			name:       "海贼王",
			args:       args{bangumiID: 975, ep: 509},
			wantEntity: &Entity{Eps: 1056},
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
	}
	db := cache.NewBolt()
	db.Open("data/bolt.db")
	anidata.Init(&anidata.Options{Cache: db})
	bangumiCache := cache.NewBolt()
	bangumiCache.Open("data/bolt_sub.db")

	store.Init(&store.InitOptions{
		Cache:        db,
		BangumiCache: bangumiCache,
	})

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

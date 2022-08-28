package bangumi

import (
	"AnimeGo/internal/cache"
	"AnimeGo/pkg/anisource"
	"fmt"
	"testing"
)

func TestBangumi_Parse(t *testing.T) {
	type args struct {
		bangumiID int
		ep        int
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *Entity
		wantEpInfo *Ep
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &Entity{Eps: 12},
			wantEpInfo: &Ep{EpID: 1109152},
			wantErr:    false,
		},
		{
			name:       "海贼王",
			args:       args{bangumiID: 975, ep: 509},
			wantEntity: &Entity{Eps: 1056},
			wantEpInfo: &Ep{EpID: 98996},
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bangumi{}
			gotEntity, gotEpInfo, err := b.Parse(tt.args.bangumiID, tt.args.ep)
			fmt.Println(gotEntity, gotEpInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEntity.Eps != tt.wantEntity.Eps {
				t.Errorf("Parse() gotEntity = %v, want %v", gotEntity, tt.wantEntity)
			}
			if gotEpInfo.EpID != tt.wantEpInfo.EpID {
				t.Errorf("Parse() gotEpInfo = %v, want %v", gotEpInfo, tt.wantEpInfo)
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
		wantEpInfo *Ep
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "联盟空军航空魔法音乐队 光辉魔女",
			args:       args{bangumiID: 253047, ep: 3},
			wantEntity: &Entity{Eps: 12},
			wantEpInfo: &Ep{EpID: 1109152},
			wantErr:    false,
		},
		{
			name:       "海贼王",
			args:       args{bangumiID: 975, ep: 509},
			wantEntity: &Entity{Eps: 1057},
			wantEpInfo: &Ep{EpID: 98996},
			wantErr:    false,
		},
	}

	db := cache.NewBolt()
	db.Open(".")
	anisource.Init(anisource.Options{Cache: db})
	b := &Bangumi{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotEntity, gotEpInfo, err := b.ParseCache(tt.args.bangumiID, tt.args.ep)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEntity.Eps != tt.wantEntity.Eps {
				t.Errorf("Parse() gotEntity = %v, want %v", gotEntity, tt.wantEntity)
			}
			if gotEpInfo.EpID != tt.wantEpInfo.EpID {
				t.Errorf("Parse() gotEpInfo = %v, want %v", gotEpInfo, tt.wantEpInfo)
			}
		})
	}
}

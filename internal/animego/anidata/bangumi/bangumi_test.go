package bangumi

import (
	"github.com/wetor/AnimeGo/test"
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

			name:       "海贼王",
			args:       args{bangumiID: 975, ep: 509},
			wantEntity: &Entity{Eps: 1057},
			wantErr:    false,
		},
	}

	test.TestInit()
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

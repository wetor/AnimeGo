package mikan

import (
	"github.com/wetor/AnimeGo/pkg/anisource"
	"github.com/wetor/AnimeGo/pkg/cache"
	"testing"
)

func TestMikan_Parse(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name          string
		args          args
		wantMikanID   int
		wantBangumiID int
		wantErr       bool
	}{
		// TODO: Add test cases.
		{
			name:          "海贼王",
			args:          args{url: "https://mikanani.me/Home/Episode/18b60d48a72c603b421468aade7fdd0868ff2f2f"},
			wantMikanID:   228,
			wantBangumiID: 975,
			wantErr:       false,
		},
		{
			name:          "欢迎来到实力至上主义的教室 第二季",
			args:          args{url: "https://mikanani.me/Home/Episode/8849c25e05d6e2623b5333bc78d3a489a9b1cc59"},
			wantMikanID:   2775,
			wantBangumiID: 371546,
			wantErr:       false,
		},
	}
	db := cache.NewBolt()
	db.Open("bolt.db")
	anisource.Init(&anisource.Options{Cache: db})
	m := &Mikan{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotMikanID, gotBangumiID, err := m.ParseCache(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMikanID != tt.wantMikanID {
				t.Errorf("Parse() gotMikanID = %v, want %v", gotMikanID, tt.wantMikanID)
			}
			if gotBangumiID != tt.wantBangumiID {
				t.Errorf("Parse() gotBangumiID = %v, want %v", gotBangumiID, tt.wantBangumiID)
			}
		})
	}
}

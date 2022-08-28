package themoviedb

import (
	"AnimeGo/data"
	"AnimeGo/internal/cache"
	"AnimeGo/pkg/anisource"
	"testing"
)

func TestThemoviedb_Parse(t1 *testing.T) {
	type fields struct {
		Key string
	}
	type args struct {
		name    string
		airDate string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantTmdbID int
		wantSeason int
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "海贼王",
			fields:     fields{Key: data.ThemoviedbKey},
			args:       args{name: "ONE PIECE", airDate: "1999-10-20"},
			wantTmdbID: 37854,
			wantSeason: 1,
			wantErr:    false,
		},
		{
			name:       "在地下城寻求邂逅是否搞错了什么 Ⅳ 新章 迷宫篇",
			fields:     fields{Key: data.ThemoviedbKey},
			args:       args{name: "ダンジョンに出会いを求めるのは間違っているだろうか Ⅳ 新章 迷宮篇", airDate: "2022-07-21"},
			wantTmdbID: 62745,
			wantSeason: 4,
			wantErr:    false,
		},
		{
			name:       "来自深渊 烈日的黄金乡",
			fields:     fields{Key: data.ThemoviedbKey},
			args:       args{name: "メイドインアビス 烈日の黄金郷", airDate: "2022-07-06"},
			wantTmdbID: 72636,
			wantSeason: 2,
			wantErr:    false,
		},
	}
	db := cache.NewBolt()
	db.Open(".")
	anisource.Init(anisource.Options{Cache: db})
	t := &Themoviedb{
		Key: data.ThemoviedbKey,
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			gotTmdbID, gotSeason, err := t.ParseCache(tt.args.name, tt.args.airDate)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTmdbID != tt.wantTmdbID {
				t1.Errorf("Parse() gotTmdbID = %v, want %v", gotTmdbID, tt.wantTmdbID)
			}
			if gotSeason != tt.wantSeason {
				t1.Errorf("Parse() gotSeason = %v, want %v", gotSeason, tt.wantSeason)
			}
		})
	}
}

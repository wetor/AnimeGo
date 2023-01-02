package themoviedb

import (
	"encoding/csv"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/key"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/request"
	"go.uber.org/zap"
	"io"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	m.Run()
	fmt.Println("end")
}

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
			fields:     fields{Key: key.ThemoviedbKey},
			args:       args{name: "ONE PIECE", airDate: "1999-10-20"},
			wantTmdbID: 37854,
			wantSeason: 1,
			wantErr:    false,
		},
		{
			name:       "在地下城寻求邂逅是否搞错了什么 Ⅳ 新章 迷宫篇",
			fields:     fields{Key: key.ThemoviedbKey},
			args:       args{name: "ダンジョンに出会いを求めるのは間違っているだろうか Ⅳ 新章 迷宮篇", airDate: "2022-07-21"},
			wantTmdbID: 62745,
			wantSeason: 4,
			wantErr:    false,
		},
		{
			name:       "来自深渊 烈日的黄金乡",
			fields:     fields{Key: key.ThemoviedbKey},
			args:       args{name: "メイドインアビス 烈日の黄金郷", airDate: "2022-07-06"},
			wantTmdbID: 72636,
			wantSeason: 2,
			wantErr:    false,
		},
		{
			name:       "OVERLORD IV",
			fields:     fields{Key: key.ThemoviedbKey},
			args:       args{name: "オーバーロードIV", airDate: "2022-07-05"},
			wantTmdbID: 64196,
			wantSeason: 4,
			wantErr:    false,
		},
		{
			name:       "福星小子",
			fields:     fields{Key: key.ThemoviedbKey},
			args:       args{name: "うる星やつら", airDate: "2022-10-14"},
			wantTmdbID: 154524,
			wantSeason: 1,
			wantErr:    false,
		},
		{
			name:       "Mairimashita! Iruma-kun 3rd Season",
			fields:     fields{Key: key.ThemoviedbKey},
			args:       args{name: "Mairimashita! Iruma-kun 3rd Season", airDate: "2022-11-14"},
			wantTmdbID: 91801,
			wantSeason: 3,
			wantErr:    false,
		},
		//
	}
	db := cache.NewBolt()
	db.Open("data/bolt.db")
	anidata.Init(&anidata.Options{Cache: db})
	t := &Themoviedb{
		Key: key.ThemoviedbKey,
	}
	request.Init(&request.InitOptions{
		Proxy: "http://127.0.0.1:7890",
	})
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			gotTmdbID, gotSeason := t.ParseCache(tt.args.name, tt.args.airDate)

			if gotTmdbID.ID != tt.wantTmdbID {
				t1.Errorf("Parse() gotTmdbID = %v, want %v", gotTmdbID, tt.wantTmdbID)
			}
			if gotSeason.Season != tt.wantSeason {
				t1.Errorf("Parse() gotSeason = %v, want %v", gotSeason, tt.wantSeason)
			}
		})
	}
}

func TestThemoviedb_ParseByFile(t1 *testing.T) {
	filename := "data/202207[20220904].csv"
	// 每个用例间隔 ms
	caseSleepMS := 100
	type args struct {
		name    string
		airDate string
	}
	type testCase struct {
		name       string
		args       args
		wantTmdbID int
		wantSeason int
		wantErr    bool
	}
	file, _ := os.Open(filename)
	defer file.Close()
	tests := make([]testCase, 0, 32)
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return
		}
		tmdbID, _ := strconv.Atoi(record[3])
		season, _ := strconv.Atoi(record[4])
		hasErr := false
		if record[5] == "true" {
			hasErr = true
		}
		tests = append(tests, testCase{
			name: record[0],
			args: args{
				name:    record[1],
				airDate: record[2],
			},
			wantTmdbID: tmdbID,
			wantSeason: season,
			wantErr:    hasErr,
		})
	}
	db := cache.NewBolt()
	db.Open("data/bolt.db")
	anidata.Init(&anidata.Options{Cache: db})
	t := &Themoviedb{
		Key: key.ThemoviedbKey,
	}
	request.Init(&request.InitOptions{
		Proxy: "http://127.0.0.1:7890",
	})
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			log.Printf("搜索：「%s」", tt.args.name)
			gotTmdbID, gotSeason := t.ParseCache(tt.args.name, tt.args.airDate)

			if gotTmdbID.ID != tt.wantTmdbID {
				t1.Errorf("Parse() gotTmdbID = %v, want %v", gotTmdbID, tt.wantTmdbID)
			}
			if gotSeason.Season != tt.wantSeason {
				t1.Errorf("Parse() gotSeason = %v, want %v", gotSeason, tt.wantSeason)
			}
			time.Sleep(time.Duration(caseSleepMS) * time.Millisecond)
		})
	}
}

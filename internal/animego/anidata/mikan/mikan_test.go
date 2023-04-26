package mikan_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "mikan"

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/test.log",
		Debug: true,
	})
	test.HookAll(testdata, nil)
	defer test.UnHook()

	db := cache.NewBolt()
	db.Open("data/bolt.db")
	anidata.Init(&anidata.Options{Cache: db})

	m.Run()

	db.Close()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestMikan_ParseCache(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name          string
		args          args
		wantMikanID   int
		wantBangumiID int
	}{
		// TODO: Add test cases.
		{
			name:          "海贼王",
			args:          args{url: "https://mikanani.me/Home/Episode/18b60d48a72c603b421468aade7fdd0868ff2f2f"},
			wantMikanID:   228,
			wantBangumiID: 975,
		},
		{
			name:          "欢迎来到实力至上主义的教室 第二季",
			args:          args{url: "https://mikanani.me/Home/Episode/8849c25e05d6e2623b5333bc78d3a489a9b1cc59"},
			wantMikanID:   2775,
			wantBangumiID: 371546,
		},
	}
	m := &mikan.Mikan{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMikanID, gotBangumiID := m.ParseCache(tt.args.url)
			if gotMikanID != tt.wantMikanID {
				t.Errorf("ParseCache() gotMikanID = %v, want %v", gotMikanID, tt.wantMikanID)
			}
			if gotBangumiID != tt.wantBangumiID {
				t.Errorf("ParseCache() gotBangumiID = %v, want %v", gotBangumiID, tt.wantBangumiID)
			}
		})
	}
}

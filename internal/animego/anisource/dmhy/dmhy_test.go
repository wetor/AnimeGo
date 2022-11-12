package dmhy

import (
	"encoding/json"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/key"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestParseDmhy(t *testing.T) {
	type args struct {
		name    string
		pubDate string
	}
	tests := []struct {
		name      string
		args      args
		wantAnime *models.AnimeEntity
	}{
		// TODO: Add test cases.
		{
			name: "圣剑传说",
			args: args{
				name:    "【MMSUB】[圣剑传说 玛娜传奇 -The Teardrop Crystal-][06][WebRip 1080p HEVC-10bit AAC][简繁内封字幕]",
				pubDate: "Sat, 12 Nov 2022 09:50:05 +0800",
			},
			wantAnime: &models.AnimeEntity{
				Season:       1,
				ThemoviedbID: 128237,
				Ep:           6,
			},
		},
		{
			name: "我家師傅沒有尾巴",
			args: args{
				name:    "[天月搬運組] 我家師傅沒有尾巴 / Uchi no Shishou wa Shippo ga Nai - 07 [1080P][簡繁日外掛] ",
				pubDate: "Sat, 12 Nov 2022 07:05:01 +0800",
			},
			wantAnime: &models.AnimeEntity{
				Season:       1,
				ThemoviedbID: 130456,
				Ep:           7,
			},
		},
		{
			name: "宠物小精灵",
			args: args{
				name:    "【枫叶字幕组】[宠物小精灵 / 宝可梦 旅途][132][简体][1080P][MP4] ",
				pubDate: "Sat, 12 Nov 2022 07:05:01 +0800",
			},
			wantAnime: &models.AnimeEntity{
				Season:       1,
				ThemoviedbID: 130456,
				Ep:           7,
			},
		},
	}
	test.TestInit()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAnime := ParseDmhy(tt.args.name, tt.args.pubDate, key.ThemoviedbKey)
			data1, _ := json.Marshal(gotAnime)
			t.Log(string(data1))
			if gotAnime.Season != tt.wantAnime.Season {
				t.Errorf("ParseDmhy().Season = %v, want %v", gotAnime.Season, tt.wantAnime.Season)
			}
			if gotAnime.Ep != tt.wantAnime.Ep {
				t.Errorf("ParseDmhy().Ep = %v, want %v", gotAnime.Ep, tt.wantAnime.Ep)
			}
			if gotAnime.ThemoviedbID != tt.wantAnime.ThemoviedbID {
				t.Errorf("ParseDmhy().ThemoviedbID = %v, want %v", gotAnime.ThemoviedbID, tt.wantAnime.ThemoviedbID)
			}
		})
	}
}

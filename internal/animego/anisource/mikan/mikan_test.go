package mikan

import (
	"AnimeGo/data"
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/cache"
	"AnimeGo/internal/models"
	"encoding/json"
	"testing"
)

func TestParseMikan(t *testing.T) {
	type args struct {
		name string
		url  string
	}
	tests := []struct {
		name      string
		args      args
		wantAnime *models.AnimeEntity
	}{
		// TODO: Add test cases.
		{
			name: "海贼王",
			args: args{
				url:  "https://mikanani.me/Home/Episode/18b60d48a72c603b421468aade7fdd0868ff2f2f",
				name: "OPFans枫雪动漫][ONE PIECE 海贼王][第1029话][1080p][周日版][MP4][简体] [299.5MB]",
			},
			wantAnime: &models.AnimeEntity{
				AnimeExtra: &models.AnimeExtra{
					MikanID:      228,
					ThemoviedbID: 37854,
				},
				ID: 975,
			},
		},
		{
			name: "欢迎来到实力至上主义的教室 第二季",
			args: args{
				url:  "https://mikanani.me/Home/Episode/8849c25e05d6e2623b5333bc78d3a489a9b1cc59",
				name: "[ANi] Classroom of the Elite S2 - 欢迎来到实力至上主义的教室 第二季 - 07 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4] [254.26 MB]",
			},
			wantAnime: &models.AnimeEntity{
				AnimeExtra: &models.AnimeExtra{
					MikanID:      2775,
					ThemoviedbID: 72517,
				},
				ID: 371546,
			},
		},
	}
	db := cache.NewBolt()
	db.Open(".")
	anisource.Init(db, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAnime := ParseMikan(tt.args.name, tt.args.url, data.ThemoviedbKey)
			data, _ := json.Marshal(gotAnime)
			t.Log(string(data))
			if gotAnime.MikanID != tt.wantAnime.MikanID {
				t.Errorf("ParseMikan().MikanID = %v, want %v", gotAnime.MikanID, tt.wantAnime.MikanID)
			}
			if gotAnime.ID != tt.wantAnime.ID {
				t.Errorf("ParseMikan().ID = %v, want %v", gotAnime.ID, tt.wantAnime.ID)
			}
			if gotAnime.ThemoviedbID != tt.wantAnime.ThemoviedbID {
				t.Errorf("ParseMikan().ThemoviedbID = %v, want %v", gotAnime.ThemoviedbID, tt.wantAnime.ThemoviedbID)
			}

		})
	}
}

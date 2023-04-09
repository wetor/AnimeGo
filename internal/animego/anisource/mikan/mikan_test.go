package mikan_test

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"sync"
	"testing"

	"github.com/brahma-adshonor/gohook"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/utils"
)

func HookGetWriter(uri string, w io.Writer) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := path.Base(uri)
	jsonData, err := os.ReadFile(path.Join("../../anidata/mikan/testdata", id+".html"))
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func HookGet(uri string, body interface{}) error {
	log.Infof("Mock HTTP GET %s", uri)
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	id := u.Query().Get("with_text_query")
	if len(id) == 0 {
		id = path.Base(u.Path)
	}
	p := path.Join("../../anidata/themoviedb/testdata", id+".json")
	if !utils.IsExist(p) {
		p = path.Join("../../anidata/bangumi/testdata", id+".json")
	}

	jsonData, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(jsonData, body)
	return nil
}

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = utils.CreateMutiDir("data")
	plugin.Init(&plugin.Options{
		Path:  "../../../../assets/plugin",
		Debug: true,
	})
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	_ = gohook.Hook(request.GetWriter, HookGetWriter, nil)
	_ = gohook.Hook(request.Get, HookGet, nil)

	b := cache.NewBolt()
	b.Open("data/bolt.db")
	anisource.Init(&anisource.Options{
		Options: &anidata.Options{
			Cache: b,
			CacheTime: map[string]int64{
				"mikan":      int64(7 * 24 * 60 * 60),
				"bangumi":    int64(3 * 24 * 60 * 60),
				"themoviedb": int64(14 * 24 * 60 * 60),
			},
		},
		TMDBFailSkip:           false,
		TMDBFailUseTitleSeason: true,
		TMDBFailUseFirstSeason: true,
	})
	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open("../../../../test/testdata/bolt_sub.bolt")
	bangumiCache.Add("bangumi_sub")
	mutex := sync.Mutex{}
	anidata.Init(&anidata.Options{
		Cache:            b,
		BangumiCache:     bangumiCache,
		BangumiCacheLock: &mutex,
	})

	request.Init(&request.Options{
		Proxy: "http://127.0.0.1:7890",
	})
	m.Run()
	b.Close()
	bangumiCache.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestMikan_Parse(t *testing.T) {
	type args struct {
		opts *models.AnimeParseOptions
		name string
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
				opts: &models.AnimeParseOptions{
					Url: "https://mikanani.me/Home/Episode/18b60d48a72c603b421468aade7fdd0868ff2f2f",
				},
				name: "OPFans枫雪动漫][ONE PIECE 海贼王][第1029话][1080p][周日版][MP4][简体] [299.5MB]",
			},
			wantAnime: &models.AnimeEntity{ID: 975, ThemoviedbID: 37854, MikanID: 228, Name: "ONE PIECE", NameCN: "海贼王", Season: 1, Ep: []*models.AnimeEpEntity{{Ep: 1029}}, Eps: 1079, AirDate: "1999-10-20"},
		},
		{
			name: "欢迎来到实力至上主义的教室 第二季",
			args: args{
				opts: &models.AnimeParseOptions{
					Url: "https://mikanani.me/Home/Episode/8849c25e05d6e2623b5333bc78d3a489a9b1cc59",
				},
				name: "[ANi] Classroom of the Elite S2 - 欢迎来到实力至上主义的教室 第二季 - 07 [1080P][Baha][WEB-DL][AAC AVC][CHT][MP4] [254.26 MB]",
			},
			wantAnime: &models.AnimeEntity{ID: 371546, ThemoviedbID: 72517, MikanID: 2775, Name: "ようこそ実力至上主義の教室へ 2nd Season", NameCN: "欢迎来到实力至上主义教室 第二季", Season: 2, Ep: []*models.AnimeEpEntity{{Ep: 7}}, Eps: 13, AirDate: "2022-07-04"},
		},
		{
			name: "想要成为影之实力者",
			args: args{
				opts: &models.AnimeParseOptions{
					Url: "https://mikanani.me/Home/Episode/dcc28079dfda415cdcdf46159aad0fa94f1a2f11",
				},
				name: "[LoliHouse] 想要成为影之实力者 / 我想成为影之强者 / Kage no Jitsuryokusha ni Naritakute! - 19 [WebRip 1080p HEVC-10bit AAC][简繁内封字幕]",
			},
			wantAnime: &models.AnimeEntity{ID: 329114, ThemoviedbID: 119495, MikanID: 2822, Name: "陰の実力者になりたくて！", NameCN: "想要成为影之实力者！", Season: 1, Ep: []*models.AnimeEpEntity{{Ep: 19}}, Eps: 20, AirDate: "2022-10-05"},
		},
	}
	m := mikan.Mikan{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := public.ParserName(tt.args.name)
			assert.NotEmpty(t, p)
			gotAnime := m.Parse(tt.args.opts)
			gotAnime.Ep = []*models.AnimeEpEntity{{Ep: p.Ep}}
			assert.Equalf(t, tt.wantAnime, gotAnime, "Parse(%v)", tt.args.opts)
		})
	}
}

package parser_test

import (
	"fmt"
	"net/url"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/anidata"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/parser"
	parserPlugin "github.com/wetor/AnimeGo/internal/animego/parser/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/torrent"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "parser"

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.MkdirAll("data", os.ModePerm)
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	torrent.Init(&torrent.Options{
		TempPath: "data",
	})
	test.HookGetWriter(testdata, nil)
	test.HookGet(testdata, func(uri string) string {
		u, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		id := u.Query().Get("with_text_query")
		if len(id) == 0 {
			id = xpath.Base(u.Path)
		}
		return id
	})
	defer test.UnHook()
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})

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
	})
	parser.Init(&parser.Options{
		TMDBFailSkip:           false,
		TMDBFailUseTitleSeason: true,
		TMDBFailUseFirstSeason: true,
	})
	bangumiCache := cache.NewBolt(true)
	bangumiCache.Open(test.GetDataPath("", "bolt_sub.bolt"))
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
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func TestManager_Parse(t *testing.T) {
	type args struct {
		opts *models.ParseOptions
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *models.AnimeEntity
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				opts: &models.ParseOptions{
					Title:      "[猎户不鸽压制] 万事屋斋藤先生转生异世界 / 斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku [01-12] [合集] [WebRip 1080p] [繁中内嵌] [H265 AAC] [2023年1月番] [4.8 GB]",
					TorrentUrl: "https://mikanani.me/Download/20230328/ac5d8d6fcc4d83cb18f18c209b66afd8e1edba86.torrent",
					MikanUrl:   "https://mikanani.me/Home/Episode/ac5d8d6fcc4d83cb18f18c209b66afd8e1edba86",
				},
			},
			wantEntity: &models.AnimeEntity{ID: 366165, ThemoviedbID: 155942, MikanID: 2922, Name: "便利屋斎藤さん、異世界に行く", NameCN: "万事屋斋藤、到异世界", Season: 1, Eps: 12, AirDate: "2023-01-08",
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpNormal, Ep: 1, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [01] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 2, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [02] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 3, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [03] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 4, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [04] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 5, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [05] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 6, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [06] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 7, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [07] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 8, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [08] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 9, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [09] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 10, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [10] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 11, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [11] [1080p] [H265 AAC] [CHT].mp4"},
					{Type: models.AnimeEpNormal, Ep: 12, Src: "[orion origin] Benriya Saitou-san, Isekai ni Iku [01-12] [WebRip] [1080p] [H265 AAC] [CHT]/[orion origin] Benriya Saitou-san, Isekai ni Iku [12] [END] [1080p] [H265 AAC] [CHT].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "ac5d8d6fcc4d83cb18f18c209b66afd8e1edba86",
					Url:  "data/ac5d8d6fcc4d83cb18f18c209b66afd8e1edba86.torrent",
				},
			},
		},
		{
			name: "2",
			args: args{
				opts: &models.ParseOptions{
					Title:      "【SW字幕组】[宠物小精灵 / 宝可梦 地平线 莉可与罗伊的旅途][01-02][简日双语字幕][2023.04.14][1080P][AVC][MP4][CHS_JP] [875.7MB]",
					TorrentUrl: "https://mikanani.me/Download/20230427/51ecf2415af99521d07595178685587e16edd926.torrent",
					MikanUrl:   "https://mikanani.me/Home/Episode/51ecf2415af99521d07595178685587e16edd926",
				},
			},
			wantEntity: &models.AnimeEntity{ID: 411247, ThemoviedbID: 220150, MikanID: 3015, Name: "ポケットモンスター", NameCN: "宝可梦 地平线", Season: 1, Eps: 22, AirDate: "2023-04-14",
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpUnknown, Ep: 0, Src: "[SWSUB][Pokemon Horizons][01-02][CHS_JP][AVC][1080P].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "51ecf2415af99521d07595178685587e16edd926",
					Url:  "data/51ecf2415af99521d07595178685587e16edd926.torrent",
				},
			},
		},
	}

	p := parserPlugin.NewParserPlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_parser.py",
	}, true)
	m := parser.NewManager(p, mikan.Mikan{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity := m.Parse(tt.args.opts)
			assert.Equalf(t, tt.wantEntity, gotEntity, "Parse(%v)", tt.args.opts)
		})
	}
}

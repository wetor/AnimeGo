package parser_test

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/parser"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/pkg/torrent"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/cache"
	pkgExceptions "github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"github.com/wetor/AnimeGo/test"
)

const testdata = "parser"

var mgr *parser.Manager

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
			id = path.Base(xpath.P(u.Path))
		}
		return id
	})
	test.Hook(torrent.LoadUri, HookLoadUri)
	defer test.UnHook()

	b := cache.NewBolt()
	b.Open("data/bolt.db")
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
		Mikan: mikan.NewMikan(&mikan.Options{
			Cache:     b,
			CacheTime: int64(7 * 24 * 60 * 60),
		}),
	})

	p := parser.NewParserPlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_parser.py",
	})
	mikanSource := &anisource.Mikan{}
	bangumiSource := &anisource.Bangumi{}

	test.HookMethod(mikanSource, "Parse", MikanParse)
	test.HookMethod(bangumiSource, "Parse", BangumiParse)

	mgr = parser.NewManager(&models.ParserOptions{
		TMDBFailSkip:           false,
		TMDBFailUseTitleSeason: true,
		TMDBFailUseFirstSeason: true,
	}, p, mikanSource, bangumiSource)

	m.Run()

	b.Close()
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
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name: "success",
			args: args{
				opts: &models.ParseOptions{
					Title:      "success",
					TorrentUrl: "success",
					MikanUrl:   "success",
				},
			},
			wantEntity: &models.AnimeEntity{ID: 366165, ThemoviedbID: 155942, MikanID: 2922, Name: "便利屋斎藤さん、異世界に行く", NameCN: "万事屋斋藤、到异世界", Season: 1, Eps: 12, AirDate: "2023-01-08",
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpNormal, Ep: 10, Src: "514/[orion origin] Benriya Saitou-san, Isekai ni Iku [10] [1080p] [H265 AAC] [CHT].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "success",
					Url:  "success",
				},
			},
		},
		{
			name: "success_bangumi",
			args: args{
				opts: &models.ParseOptions{
					Title:      "success",
					TorrentUrl: "success",
					MikanUrl:   "success",
					BangumiID:  366165,
				},
			},
			wantEntity: &models.AnimeEntity{ID: 366165, ThemoviedbID: 155942, MikanID: 2922, Name: "便利屋斎藤さん、異世界に行く", NameCN: "万事屋斋藤、到异世界", Season: 1, Eps: 12, AirDate: "2023-01-08",
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpNormal, Ep: 10, Src: "514/[orion origin] Benriya Saitou-san, Isekai ni Iku [10] [1080p] [H265 AAC] [CHT].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "success",
					Url:  "success",
				},
			},
		},
		{
			name: "ep_unknown",
			args: args{
				opts: &models.ParseOptions{
					Title:      "ep_unknown",
					TorrentUrl: "ep_unknown",
					MikanUrl:   "ep_unknown",
				},
			},
			wantEntity: &models.AnimeEntity{ID: 411247, ThemoviedbID: 220150, MikanID: 3015, Name: "ポケットモンスター", NameCN: "宝可梦 地平线", Season: 1, Eps: 22, AirDate: "2023-04-14",
				Flag: models.AnimeFlagEpParseFailed,
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpUnknown, Ep: 0, Src: "[SWSUB][Pokemon Horizons][01-02][CHS_JP][AVC][1080P].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "ep_unknown",
					Url:  "ep_unknown",
				},
			},
		},
		{
			name: "err_anisource_parse",
			args: args{
				opts: &models.ParseOptions{
					Title:      "err_anisource_parse",
					TorrentUrl: "err_anisource_parse",
					MikanUrl:   "err_anisource_parse",
				},
			},
			wantErr:    &exceptions.ErrMikanParseHTML{},
			wantErrStr: "解析anisource失败，结束此流程: 解析Mikan信息失败: 解析 Input 失败，解析网页错误",
		},
		{
			name: "err_torrent",
			args: args{
				opts: &models.ParseOptions{
					Title:      "err_torrent",
					TorrentUrl: "err_torrent",
					MikanUrl:   "err_torrent",
				},
			},
			wantErr:    &pkgExceptions.ErrTorrentUrl{},
			wantErrStr: "解析torrent失败，结束此流程: 无法识别Torrent Url: err_torrent",
		},
		{
			name: "err_season",
			args: args{
				opts: &models.ParseOptions{
					Title:      "err_season",
					TorrentUrl: "err_season",
					MikanUrl:   "err_season",
				},
			},
			wantEntity: &models.AnimeEntity{ID: 411247, ThemoviedbID: 220150, MikanID: 3015, Name: "ポケットモンスター", NameCN: "宝可梦 地平线", Season: 1, Eps: 22, AirDate: "2023-04-14",
				Flag: models.AnimeFlagEpParseFailed,
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpUnknown, Ep: 0, Src: "[SWSUB][Pokemon Horizons][01-02][CHS_JP][AVC][1080P].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "err_season",
					Url:  "err_season",
				},
			},
		},
		{
			name: "err_season_use_title",
			args: args{
				opts: &models.ParseOptions{
					Title:      "err_season_use_title",
					TorrentUrl: "err_season_use_title",
					MikanUrl:   "err_season_use_title",
				},
			},
			wantEntity: &models.AnimeEntity{ID: 411247, ThemoviedbID: 220150, MikanID: 3015, Name: "ポケットモンスター", NameCN: "宝可梦 地平线", Season: 2, Eps: 22, AirDate: "2023-04-14",
				Ep: []*models.AnimeEpEntity{
					{Type: models.AnimeEpNormal, Ep: 1, Src: "[SWSUB][Pokemon Horizons][第二季][01][CHS_JP][AVC][1080P].mp4"},
				},
				Torrent: &models.AnimeTorrent{
					Hash: "err_season_use_title",
					Url:  "err_season_use_title",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := mgr.Parse(tt.args.opts)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantEntity, gotEntity, "Parse(%v)", tt.args.opts)
			}
		})
	}
}

func TestManager_Parse_Failed(t *testing.T) {
	type args struct {
		opts *models.ParseOptions
	}
	tests := []struct {
		name       string
		args       args
		wantEntity *models.AnimeEntity
		wantErr    error
		wantErrStr string
	}{
		// TODO: Add test cases.
		{
			name: "err_season_failed",
			args: args{
				opts: &models.ParseOptions{
					Title:      "err_season_failed",
					TorrentUrl: "err_season_failed",
					MikanUrl:   "err_season_failed",
				},
			},
			wantErr:    &exceptions.ErrParseFailed{},
			wantErrStr: "解析季度失败，未设置默认值，结束此流程",
		},
	}
	mgr.TMDBFailSkip = true
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntity, err := mgr.Parse(tt.args.opts)
			if tt.wantErr != nil {
				assert.IsType(t, tt.wantErr, errors.Cause(err))
				assert.EqualError(t, err, tt.wantErrStr)
			} else {
				assert.NoError(t, err)
				assert.Equalf(t, tt.wantEntity, gotEntity, "Parse(%v)", tt.args.opts)
			}
		})
	}
}

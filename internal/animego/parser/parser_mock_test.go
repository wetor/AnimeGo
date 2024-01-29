package parser_test

import (
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/pkg/torrent"
	pkgExceptions "github.com/wetor/AnimeGo/pkg/exceptions"
)

func MikanParse(mikan *anisource.Mikan, opt *models.AnimeParseOptions) (*models.AnimeEntity, error) {
	return parse(opt)
}

func BangumiParse(bgm *anisource.Bangumi, opt *models.AnimeParseOptions) (*models.AnimeEntity, error) {
	return parse(opt)
}

func parse(opt *models.AnimeParseOptions) (*models.AnimeEntity, error) {
	bangumiID := 0
	mikanUrl, ok := opt.Input.(string)
	if !ok {
		bangumiID = opt.Input.(int)
	}

	if bangumiID > 0 {
		switch bangumiID {
		case 366165:
			return &models.AnimeEntity{
				ID:           366165,
				ThemoviedbID: 155942,
				MikanID:      2922,
				Name:         "便利屋斎藤さん、異世界に行く",
				NameCN:       "万事屋斋藤、到异世界",
				Season:       1,
				Eps:          12,
				AirDate:      "2023-01-08",
			}, nil
		}
		return nil, nil
	}

	switch mikanUrl {
	case "success":
		return &models.AnimeEntity{
			ID:           366165,
			ThemoviedbID: 155942,
			MikanID:      2922,
			Name:         "便利屋斎藤さん、異世界に行く",
			NameCN:       "万事屋斋藤、到异世界",
			Season:       1,
			Eps:          12,
			AirDate:      "2023-01-08",
		}, nil
	case "ep_unknown", "err_torrent":
		return &models.AnimeEntity{
			ID:           411247,
			ThemoviedbID: 220150,
			MikanID:      3015,
			Name:         "ポケットモンスター",
			NameCN:       "宝可梦 地平线",
			Season:       1,
			Eps:          22,
			AirDate:      "2023-04-14",
		}, nil
	case "err_anisource_parse":
		return nil, errors.Wrap(&exceptions.ErrMikanParseHTML{Message: "Input"},
			"解析Mikan信息失败")
	case "err_season", "err_season_use_title", "err_season_failed":
		return &models.AnimeEntity{
			ID:           411247,
			ThemoviedbID: 220150,
			MikanID:      3015,
			Name:         "ポケットモンスター",
			NameCN:       "宝可梦 地平线",
			Season:       0,
			Eps:          22,
			AirDate:      "2023-04-14",
		}, nil
	}
	return nil, nil
}

func HookLoadUri(uri string) (t *torrent.Torrent, err error) {
	switch uri {
	case "success":
		return &torrent.Torrent{
			Type:   "torrent",
			Url:    "success",
			Hash:   "success",
			Name:   "success",
			Length: 5185368669,
			Files: []*torrent.File{
				{
					Name:   "[orion origin] Benriya Saitou-san, Isekai ni Iku [10] [1080p] [H265 AAC] [CHT].mp4",
					Dir:    "514",
					Length: 1919,
				},
			},
		}, nil
	case "ep_unknown":
		return &torrent.Torrent{
			Type:   "torrent",
			Url:    "ep_unknown",
			Hash:   "ep_unknown",
			Name:   "ep_unknown",
			Length: 5185368669,
			Files: []*torrent.File{
				{
					Name:   "[SWSUB][Pokemon Horizons][01-02][CHS_JP][AVC][1080P].mp4",
					Dir:    "",
					Length: 1919,
				},
			},
		}, nil
	case "err_torrent":
		return nil, &pkgExceptions.ErrTorrentUrl{Url: uri}
	case "err_season":
		return &torrent.Torrent{
			Type:   "torrent",
			Url:    "err_season",
			Hash:   "err_season",
			Name:   "err_season",
			Length: 5185368669,
			Files: []*torrent.File{
				{
					Name:   "[SWSUB][Pokemon Horizons][01-02][CHS_JP][AVC][1080P].mp4",
					Dir:    "",
					Length: 1919,
				},
			},
		}, nil
	case "err_season_use_title", "err_season_failed":
		return &torrent.Torrent{
			Type:   "torrent",
			Url:    "err_season_use_title",
			Hash:   "err_season_use_title",
			Name:   "err_season_use_title",
			Length: 5185368669,
			Files: []*torrent.File{
				{
					Name:   "[SWSUB][Pokemon Horizons][第二季][01][CHS_JP][AVC][1080P].mp4",
					Dir:    "",
					Length: 1919,
				},
			},
		}, nil
	}
	return nil, nil
}

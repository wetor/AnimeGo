package bangumi

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Bangumi struct {
	ThemoviedbKey string
}

func NewBangumiSource(tmdbKey string) api.AniSource {
	return Bangumi{
		ThemoviedbKey: tmdbKey,
	}
}

func (m Bangumi) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity, err error) {
	var bgmID int
	var mikanID int
	var bgmEntity = &bangumi.Entity{}
	var season int
	var tmdbID int

	switch input := opts.Input.(type) {
	case models.MikanEntity:
		bgmID = input.BangumiID
		mikanID = input.MikanID
	case int:
		bgmID = input
		mikanID = 0
	}
	// ------------------- 获取bangumi信息 -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideBangumi() {
		log.Debugf("[AniSource] 解析Bangumi，%d", bgmID)
		entity, err := anisource.Bangumi().GetCache(bgmID, nil)
		if err != nil {
			return nil, err
		}
		bgmEntity = entity.(*bangumi.Entity)
	} else {
		bgmEntity.Name = opts.AnimeParseOverride.Name
		bgmEntity.NameCN = opts.AnimeParseOverride.NameCN
		bgmEntity.AirDate = opts.AnimeParseOverride.AirDate
		bgmEntity.Eps = opts.AnimeParseOverride.Eps
	}

	// ------------------- 获取tmdb信息(季度信息) -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideThemoviedb() {
		log.Debugf("[AniSource] 解析Themoviedb，%s, %s", bgmEntity.Name, bgmEntity.AirDate)
		t := anisource.Themoviedb(m.ThemoviedbKey)
		id, err := t.SearchCache(bgmEntity.Name, nil)
		if err != nil {
			return nil, err
		}
		tmdbID = id
		entity, err := t.GetCache(id, bgmEntity.AirDate)
		if err != nil {
			log.Warnf("[AniSource] 解析Themoviedb获取番剧季度信息失败")
		} else {
			season = entity.(*themoviedb.SeasonInfo).Season
		}
	} else {
		tmdbID = opts.AnimeParseOverride.ThemoviedbID
		season = opts.AnimeParseOverride.Season
	}

	anime = &models.AnimeEntity{
		ID:           bgmID,
		ThemoviedbID: tmdbID,
		MikanID:      mikanID,
		Name:         bgmEntity.Name,
		NameCN:       bgmEntity.NameCN,
		Season:       season,
		Eps:          bgmEntity.Eps,
		AirDate:      bgmEntity.AirDate,
	}
	anime.Default()
	log.Infof("[AniSource] 获取「%s」信息成功！", anime.NameCN)
	return anime, nil
}

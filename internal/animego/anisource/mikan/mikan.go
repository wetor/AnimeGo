package mikan

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Mikan struct {
	ThemoviedbKey string
}

func (m Mikan) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity, err error) {
	var mikanEntity = &mikan.Entity{}
	var bgmEntity = &bangumi.Entity{}
	var season int
	var tmdbID int

	// ------------------- 获取bangumiID -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideMikan() {
		log.Debugf("步骤1，解析Mikan，%s", opts.MikanUrl)
		entity, err := anisource.Mikan().ParseCache(opts.MikanUrl)
		if err != nil {
			return nil, err
		}
		mikanEntity = entity.(*mikan.Entity)
	} else {
		mikanEntity.MikanID = opts.AnimeParseOverride.MikanID
		mikanEntity.BangumiID = opts.AnimeParseOverride.BangumiID
	}
	// ------------------- 获取bangumi信息 -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideBangumi() {
		log.Debugf("步骤2，解析Bangumi，%d", mikanEntity.BangumiID)
		entity, err := anisource.Bangumi().GetCache(mikanEntity.BangumiID, nil)
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
		log.Debugf("步骤3，解析Themoviedb，%s, %s", bgmEntity.Name, bgmEntity.AirDate)
		t := anisource.Themoviedb(m.ThemoviedbKey)
		id, err := t.SearchCache(bgmEntity.Name)
		if err != nil {
			return nil, err
		}
		tmdbID = id
		entity, err := t.GetCache(id, bgmEntity.AirDate)
		if err != nil {
			log.Warnf("解析Themoviedb获取番剧季度信息失败")
		} else {
			season = entity.(*themoviedb.SeasonInfo).Season
		}
	} else {
		tmdbID = opts.AnimeParseOverride.ThemoviedbID
		season = opts.AnimeParseOverride.Season
	}

	anime = &models.AnimeEntity{
		ID:           mikanEntity.BangumiID,
		ThemoviedbID: tmdbID,
		MikanID:      mikanEntity.MikanID,
		Name:         bgmEntity.Name,
		NameCN:       bgmEntity.NameCN,
		Season:       season,
		Eps:          bgmEntity.Eps,
		AirDate:      bgmEntity.AirDate,
	}
	anime.Default()
	log.Infof("获取「%s」信息成功！", anime.NameCN)
	return anime, nil
}

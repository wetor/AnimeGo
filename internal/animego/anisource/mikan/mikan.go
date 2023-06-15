package mikan

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anidata/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

type Mikan struct {
	ThemoviedbKey string
}

func (m Mikan) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity) {
	var err error
	var season int
	var mikanEntity = &mikan.Entity{}
	var entity = &bangumi.Entity{}
	var tmdbEntity = &themoviedb.Entity{}
	var tmdbSeason = &themoviedb.SeasonInfo{}

	// ------------------- 获取bangumiID -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideMikan() {
		log.Debugf("步骤1，解析Mikan，%s", opts.MikanUrl)
		try.This(func() {
			mikanEntity = anisource.Mikan().ParseCache(opts.MikanUrl).(*mikan.Entity)
		}).Catch(func(e try.E) {
			err = e.(error)
		})
		if err != nil {
			log.Warnf("解析Mikan获取bangumi id失败")
			log.Debugf("", err)
			return nil
		}
	} else {
		mikanEntity.MikanID = opts.AnimeParseOverride.MikanID
		mikanEntity.BangumiID = opts.AnimeParseOverride.BangumiID
	}
	// ------------------- 获取bangumi信息 -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideBangumi() {
		log.Debugf("步骤2，解析Bangumi，%d", mikanEntity.BangumiID)
		try.This(func() {
			entity = anisource.Bangumi().GetCache(mikanEntity.BangumiID, nil).(*bangumi.Entity)
		}).Catch(func(e try.E) {
			err = e.(error)
		})
		if err != nil {
			log.Warnf("解析bangumi获取番剧信息失败")
			log.Debugf("", err)
			return nil
		}
	} else {
		entity.Name = opts.AnimeParseOverride.Name
		entity.NameCN = opts.AnimeParseOverride.NameCN
		entity.AirDate = opts.AnimeParseOverride.AirDate
		entity.Eps = opts.AnimeParseOverride.Eps
	}

	// ------------------- 获取tmdb信息(季度信息) -------------------
	if opts.AnimeParseOverride == nil || !opts.OverrideThemoviedb() {
		log.Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
		try.This(func() {
			t := anisource.Themoviedb(m.ThemoviedbKey)
			id := t.SearchCache(entity.Name)
			tmdbEntity.ID = id
			tmdbSeason = t.GetCache(id, entity.AirDate).(*themoviedb.SeasonInfo)
			season = tmdbSeason.Season
		}).Catch(func(e try.E) {
			err = e.(error)
		})
		if err != nil {
			log.Warnf("解析Themoviedb获取番剧季度信息失败")
			log.Debugf("", err)
		}
	} else {
		tmdbEntity.ID = opts.AnimeParseOverride.ThemoviedbID
		season = opts.AnimeParseOverride.Season
	}

	anime = &models.AnimeEntity{
		ID:           mikanEntity.BangumiID,
		ThemoviedbID: tmdbEntity.ID,
		MikanID:      mikanEntity.MikanID,
		Name:         entity.Name,
		NameCN:       entity.NameCN,
		Season:       season,
		Eps:          entity.Eps,
		AirDate:      entity.AirDate,
	}
	anime.Default()
	log.Infof("获取「%s」信息成功！", anime.NameCN)
	return anime
}

package mikan

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
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
	var mikanID, bangumiID, season int
	var entity = &bangumi.Entity{}
	var tmdbEntity = &themoviedb.Entity{}
	var tmdbSeason = &themoviedb.SeasonInfo{}

	// ------------------- 获取bangumiID -------------------
	log.Debugf("步骤1，解析Mikan，%s", opts.Url)
	try.This(func() {
		mikanID, bangumiID = anisource.Mikan().ParseCache(opts.Url)
	}).Catch(func(e try.E) {
		err = e.(error)
	})
	if err != nil {
		log.Warnf("解析Mikan获取bangumi id失败")
		log.Debugf("", err)
		return nil
	}

	// ------------------- 获取bangumi信息 -------------------
	log.Debugf("步骤2，解析Bangumi，%d", bangumiID)
	try.This(func() {
		entity = anisource.Bangumi().ParseCache(bangumiID)
	}).Catch(func(e try.E) {
		err = e.(error)
	})
	if err != nil {
		log.Warnf("解析bangumi获取番剧信息失败失败")
		log.Debugf("", err)
		return nil
	}

	// ------------------- 获取tmdb信息(季度信息) -------------------
	log.Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
	try.This(func() {
		tmdbEntity, tmdbSeason = anisource.Themoviedb(m.ThemoviedbKey).ParseCache(entity.Name, entity.AirDate)
		season = tmdbSeason.Season
	}).Catch(func(e try.E) {
		err = e.(error)
	})
	if err != nil {
		log.Debugf("", err)
		return nil
	}
	anime = &models.AnimeEntity{
		ID:           entity.ID,
		ThemoviedbID: tmdbEntity.ID,
		MikanID:      mikanID,
		Name:         entity.Name,
		NameCN:       entity.NameCN,
		Season:       season,
		Eps:          entity.Eps,
		AirDate:      entity.AirDate,
	}
	anime.Default()
	log.Infof("获取「%s」信息成功！原名「%s」", anime.NameCN, anime.Name)
	return anime
}

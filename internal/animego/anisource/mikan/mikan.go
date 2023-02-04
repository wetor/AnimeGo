package mikan

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

type Mikan struct {
	ThemoviedbKey string
}

func (m Mikan) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity) {
	defer func() {
		// 利用panic结束函数
		recover()
	}()

	log.Infof("获取「%s」信息开始...", opts.Name)
	// ------------------- 解析文件名获取ep -------------------
	if opts.Parsed == nil || opts.Parsed.Ep == 0 {
		opts.Parsed = public.ParserName(opts.Name)
		if opts.Parsed.Ep == 0 {
			log.Warnf("解析ep信息失败，结束此流程")
			return nil
		}
	}
	var mikanID, bangumiID, season int
	var entity = &bangumi.Entity{}
	var tmdbEntity = &themoviedb.Entity{}
	var tmdbSeason = &themoviedb.SeasonInfo{}

	// ------------------- 获取mikanID -------------------
	log.Debugf("步骤1，解析Mikan，%s", opts.Url)
	try.This(func() {
		mikanID, bangumiID = anisource.Mikan().ParseCache(opts.Url)
	}).Catch(func(err try.E) {
		log.Warnf("解析Mikan获取bangumi id失败，结束此流程")
		log.Debugf("", err)
	})

	// ------------------- 获取bangumi信息 -------------------
	log.Debugf("步骤2，解析Bangumi，%d, %d", bangumiID, opts.Parsed.Ep)
	try.This(func() {
		entity = anisource.Bangumi().ParseCache(bangumiID)
	}).Catch(func(err try.E) {
		log.Warnf("解析bangumi获取番剧信息失败失败，结束此流程")
		log.Debugf("", err)
		panic("return")
	})

	// ------------------- 获取tmdb信息(季度信息) -------------------
	season = 1
	log.Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
	try.This(func() {
		tmdbEntity, tmdbSeason = anisource.Themoviedb(m.ThemoviedbKey).ParseCache(entity.Name, entity.AirDate)
		season = tmdbSeason.Season
	}).Catch(func(err try.E) {
		if anisource.TMDBFailSkip {
			log.Warnf("无法获取准确的季度信息，结束此流程")
			log.Debugf("", err)
			panic("return")
		} else if anisource.TMDBFailUseTitleSeason && opts.Parsed.Season != 0 {
			season = opts.Parsed.Season
			log.Warnf("使用标题解析季度信息：第%d季", opts.Parsed.Season)
		}
		if season == 0 {
			if anisource.TMDBFailUseFirstSeason {
				season = 1
				log.Warnf("无法获取准确季度信息，默认：第%d季", season)
			} else {
				log.Warnf("无法获取准确的季度信息，结束此流程")
				panic("return")
			}
		}
	})

	anime = &models.AnimeEntity{
		ID:           entity.ID,
		ThemoviedbID: tmdbEntity.ID,
		MikanID:      mikanID,
		Name:         entity.Name,
		NameCN:       entity.NameCN,
		Season:       season,
		Ep:           opts.Parsed.Ep,
		Eps:          entity.Eps,
		AirDate:      entity.AirDate,
	}
	log.Infof("获取「%s」信息成功！原名「%s」", anime.FullName(), anime.Name)
	return anime
}

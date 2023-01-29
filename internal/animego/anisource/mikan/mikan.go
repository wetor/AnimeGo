package mikan

import (
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/errors"
)

type Mikan struct {
	ThemoviedbKey string
}

func (m Mikan) Parse(opts *models.AnimeParseOptions) (anime *models.AnimeEntity) {
	errMsg := ""
	defer errors.HandleAniError(func(err *errors.AniError) {
		zap.S().Debug(err)
		zap.S().Warn(errMsg)
	})

	zap.S().Infof("获取「%s」信息开始...", opts.Name)
	// ------------------- 解析文件名获取ep -------------------
	if opts.Parsed == nil || opts.Parsed.Ep == 0 {
		opts.Parsed = public.ParserName(opts.Name)
		if opts.Parsed.Ep == 0 {
			zap.S().Warn("解析ep信息失败，结束此流程")
			return nil
		}
	}
	// ------------------- 获取mikanID -------------------
	zap.S().Debugf("步骤1，解析Mikan，%s", opts.Url)
	errMsg = "解析Mikan获取bangumi id失败，结束此流程"
	mikanID, bangumiID := anisource.Mikan().ParseCache(opts.Url)

	// ------------------- 获取bangumi信息 -------------------
	zap.S().Debugf("步骤2，解析Bangumi，%d, %d", bangumiID, opts.Parsed.Ep)
	errMsg = "解析bangumi获取番剧信息失败失败，结束此流程"
	entity := anisource.Bangumi().ParseCache(bangumiID)

	// ------------------- 获取tmdb信息(季度信息) -------------------
	zap.S().Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
	tmdbEntity := &themoviedb.Entity{}
	tmdbSeason := &themoviedb.SeasonInfo{}
	season := 1
	func() {
		defer errors.HandleAniError(func(err *errors.AniError) {
			zap.S().Debug(err)
			if anisource.TMDBFailSkip {
				zap.S().Warn("无法获取准确的季度信息，结束此流程")
				return
			}
			if anisource.TMDBFailUseTitleSeason && opts.Parsed.Season != 0 {
				season = opts.Parsed.Season
				zap.S().Debugf("使用标题解析季度信息：第%d季", opts.Parsed.Season)
			}
			if season == 0 {
				if anisource.TMDBFailUseFirstSeason {
					season = 1
					zap.S().Debugf("无法获取准确季度信息，默认：第%d季", season)
				} else {
					zap.S().Warn("无法获取准确的季度信息，结束此流程")
					return
				}
			}
		})
		tmdbEntity, tmdbSeason = anisource.Themoviedb(m.ThemoviedbKey).ParseCache(entity.Name, entity.AirDate)
		season = tmdbSeason.Season
	}()

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
	zap.S().Infof("获取「%s」信息成功！原名「%s」", anime.FullName(), anime.Name)
	return anime
}

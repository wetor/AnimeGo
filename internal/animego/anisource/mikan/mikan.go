package mikan

import (
	"github.com/wetor/AnimeGo/internal/animego/anidata/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/third_party/poketto"
	"go.uber.org/zap"
)

func ParseMikan(name, url, tmdbKey string) (anime *models.AnimeEntity) {
	errMsg := ""
	defer errors.HandleAniError(func(err *errors.AniError) {
		zap.S().Debug(err)
		zap.S().Warn(errMsg)
	})

	zap.S().Infof("获取「%s」信息开始...", name)
	// ------------------- 解析文件名获取ep -------------------
	match := poketto.NewEpisode(name)
	match.TryParse()
	if match.ParseErr == poketto.CannotParseEpErr {
		zap.S().Warn("解析ep信息失败，结束此流程")
		return nil
	}
	// ------------------- 获取mikanID -------------------
	zap.S().Debugf("步骤1，解析Mikan，%s", url)
	errMsg = "解析Mikan获取bangumi id失败，结束此流程"
	mikanID, bangumiID := anisource.Mikan().ParseCache(url)

	// ------------------- 获取bangumi信息 -------------------
	zap.S().Debugf("步骤2，解析Bangumi，%d, %d", bangumiID, match.Ep)
	errMsg = "解析bangumi获取番剧信息失败失败，结束此流程"
	entity, epInfo := anisource.Bangumi().ParseCache(bangumiID, match.Ep)

	// ------------------- 获取tmdb信息(季度信息) -------------------
	zap.S().Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
	tmdbEntity := &themoviedb.Entity{}
	tmdbSeason := &themoviedb.SeasonInfo{}
	season := 1
	func() {
		defer errors.HandleAniError(func(err *errors.AniError) {
			zap.S().Debug(err)
			if store.Config.Default.TMDBFailSkip {
				zap.S().Warn("无法获取准确的季度信息，结束此流程")
				return
			}
			if store.Config.Default.TMDBFailUseTitleSeason && match.Season != 0 {
				season = match.Season
				zap.S().Debugf("使用标题解析季度信息：第%d季", match.Season)
			}
			if season == 0 {
				if store.Config.Default.TMDBFailUseFirstSeason {
					season = 1
					zap.S().Debugf("无法获取准确季度信息，默认：第%d季", season)
				} else {
					zap.S().Warn("无法获取准确的季度信息，结束此流程")
					return
				}
			}
		})
		tmdbEntity, tmdbSeason = anisource.Themoviedb(tmdbKey).ParseCache(entity.Name, entity.AirDate)
		season = tmdbSeason.Season
	}()

	anime = &models.AnimeEntity{
		ID:           entity.ID,
		ThemoviedbID: tmdbEntity.ID,
		MikanID:      mikanID,
		Name:         entity.Name,
		NameCN:       entity.NameCN,
		Season:       season,
		Ep:           epInfo.Ep,
		EpID:         epInfo.ID,
		Eps:          entity.Eps,
		AirDate:      entity.AirDate,
		Date:         epInfo.AirDate,
	}
	zap.S().Infof("获取「%s」信息成功！原名「%s」", anime.FullName(), anime.Name)
	return anime
}

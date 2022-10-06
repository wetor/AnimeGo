package mikan

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/models"
	"AnimeGo/internal/store"
	"AnimeGo/third_party/poketto"
	"go.uber.org/zap"
)

func ParseMikan(name, url, tmdbKey string) (anime *models.AnimeEntity) {
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
	mikanID, bangumiID, err := anisource.Mikan().ParseCache(url)
	if err != nil {
		zap.S().Debug(err)
		zap.S().Warn("解析Mikan获取bangumi id失败，结束此流程")
		return nil
	}

	// ------------------- 获取bangumi信息 -------------------
	zap.S().Debugf("步骤2，解析Bangumi，%d, %d", bangumiID, match.Ep)
	entity, epInfo, err := anisource.Bangumi().ParseCache(bangumiID, match.Ep)
	if err != nil {
		zap.S().Debug(err)
		zap.S().Warn("解析bangumi获取番剧信息失败失败，结束此流程")
		return nil
	}
	// ------------------- 获取tmdb信息(季度信息) -------------------
	zap.S().Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
	tmdbID, season, err := anisource.Themoviedb(tmdbKey).ParseCache(entity.Name, entity.AirDate)
	if err != nil {
		zap.S().Debug(err)
		if store.Config.Default.TMDBFailSkip {
			zap.S().Warn("无法获取准确的季度信息，结束此流程")
			return nil
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
				return nil
			}
		}
	}

	anime = &models.AnimeEntity{
		ID:      entity.ID,
		Name:    entity.Name,
		NameCN:  entity.NameCN,
		AirDate: entity.AirDate,
		Eps:     entity.Eps,
		AnimeSeason: &models.AnimeSeason{
			Season: season,
		},
		AnimeEp: &models.AnimeEp{
			Ep:       epInfo.Ep,
			Date:     epInfo.Date,
			Duration: epInfo.Duration,
			EpDesc:   epInfo.EpDesc,
			EpName:   epInfo.EpName,
			EpNameCN: epInfo.EpNameCN,
			EpID:     epInfo.EpID,
		},
		AnimeExtra: &models.AnimeExtra{
			MikanID:      mikanID,
			MikanUrl:     url,
			ThemoviedbID: tmdbID,
		},
	}
	zap.S().Infof("获取「%s」信息成功！原名「%s」", anime.FullName(), anime.Name)
	return anime
}

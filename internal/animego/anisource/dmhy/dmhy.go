package dmhy

import (
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/third_party/poketto"
	"go.uber.org/zap"
)

func ParseDmhy(name, pubDate, tmdbKey string) (anime *models.AnimeEntity) {
	zap.S().Infof("获取「%s」信息开始...", name)
	// ------------------- 解析文件名获取ep -------------------
	match := poketto.NewEpisode(name)
	match.TryParse()
	if match.ParseErr == poketto.CannotParseEpErr {
		zap.S().Warn("解析ep信息失败，结束此流程")
		return nil
	}
	if len(match.Name) == 0 {
		zap.S().Warn("解析番剧名信息失败，结束此流程")
		return nil
	}

	pubDate = utils.UTCToTimeStr(pubDate)
	// ------------------- 获取tmdb信息(季度信息) -------------------
	zap.S().Debugf("步骤3，解析Themoviedb，%s", match.Name)
	tmdbEntity, tmdbSeason, err := anisource.Themoviedb(tmdbKey).ParseCache(match.Name, pubDate)
	season := 1
	if err != nil {
		zap.S().Debug(err)
		if tmdbEntity == nil {
			zap.S().Warn("解析失败，结束此流程")
			return nil
		}
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
	} else {
		season = tmdbSeason.Season
	}

	anime = &models.AnimeEntity{
		ID:           0,
		ThemoviedbID: tmdbEntity.ID,
		MikanID:      0,
		Name:         tmdbEntity.Name,
		NameCN:       tmdbEntity.NameCN,
		Season:       season,
		Ep:           match.Ep,
		EpID:         0,
		Date:         "",
	}
	if tmdbSeason != nil {
		anime.Eps = tmdbSeason.Eps
		anime.AirDate = tmdbSeason.AirDate
	}
	zap.S().Infof("获取「%s」信息成功！原名「%s」", anime.FullName(), anime.Name)
	return anime
}

package mikan

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/animego/parser"
	"AnimeGo/internal/models"

	"go.uber.org/zap"
)

func ParseMikan(name, url, tmdbKey string) (anime *models.AnimeEntity) {
	zap.S().Infof("获取「%s」信息开始...", name)
	// ------------------- 解析文件名获取ep -------------------
	match, err := parser.ParseTitle(name)
	if err != nil {
		zap.S().Warn("解析ep信息失败，结束此流程")
		return nil
	}
	// ------------------- 获取mikanID -------------------
	zap.S().Debugf("步骤1，解析Mikan，%s", url)
	mikanID, bangumiID, err := anisource.Mikan().ParseCache(url)
	if err != nil {
		zap.S().Warn(err)
		zap.S().Warn("结束此流程")
		return nil
	}

	// ------------------- 获取bangumi信息 -------------------
	zap.S().Debugf("步骤2，解析Bangumi，%d, %d", bangumiID, match.Ep)
	entity, epInfo, err := anisource.Bangumi().ParseCache(bangumiID, match.Ep)
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	// ------------------- 获取tmdb信息(季度信息) -------------------
	zap.S().Debugf("步骤3，解析Themoviedb，%s, %s", entity.Name, entity.AirDate)
	tmdbID, season, err := anisource.Themoviedb(tmdbKey).ParseCache(entity.Name, entity.AirDate)
	if err != nil {
		zap.S().Warn(err)
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

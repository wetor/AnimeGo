package parser

import (
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/torrent"
)

const DefaultSeason = 1

type Manager struct {
	parser    api.ParserPlugin
	anisource api.AniSource
}

func NewManager(parser api.ParserPlugin, anisource api.AniSource) *Manager {
	return &Manager{
		parser:    parser,
		anisource: anisource,
	}
}

func (m *Manager) Parse(opts *models.ParseOptions) (entity *models.AnimeEntity) {

	// ------------------- 获取mikan信息（bangumi id） -------------------
	entity = m.anisource.Parse(&models.AnimeParseOptions{
		Url: opts.MikanUrl,
	})
	if entity == nil {
		log.Warnf("结束此流程")
		return
	}
	// ------------------- 获取并解析torrent信息 -------------------
	torrentInfo, err := torrent.LoadUri(opts.TorrentUrl)
	if err != nil {
		log.Debugf("", err)
		log.Warnf("解析torrent失败，结束此流程")
		return
	}
	entity.Ep = make([]*models.AnimeEpEntity, len(torrentInfo.Files))

	entity.Torrent = &models.AnimeTorrent{
		Hash: torrentInfo.Hash,
		Url:  torrentInfo.Url,
	}
	for i, t := range torrentInfo.Files {
		file := m.parser.Parse(t.Name)
		entity.Ep[i] = &models.AnimeEpEntity{
			Ep:  file.Ep,
			Src: t.Path(),
		}
	}
	// ------------------- 解析标题获取季度信息 -------------------
	parsed := m.parser.Parse(opts.Title)

	// 优先tmdb解析的season，如果为0则根据配置使用标题解析season
	if entity.Season == 0 {
		entity.Season = m.defaultSeason(parsed.Season)
	}

	return entity
}

func (m *Manager) defaultSeason(season int) (result int) {
	if !TMDBFailSkip {
		if TMDBFailUseTitleSeason && season != 0 {
			result = season
			log.Warnf("使用标题解析季度信息：第%d季", result)
			return
		}
		if TMDBFailUseFirstSeason {
			result = DefaultSeason
			log.Warnf("无法获取准确季度信息，默认：第%d季", result)
			return
		}
	}
	log.Warnf("无法获取准确的季度信息，结束此流程")
	return
}

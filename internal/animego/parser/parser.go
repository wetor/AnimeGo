package parser

import (
	"fmt"

	"github.com/wetor/AnimeGo/internal/animego/parser/utils"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/json"
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
		return nil
	}
	// ------------------- 获取并解析torrent信息 -------------------
	torrentInfo, err := torrent.LoadUri(opts.TorrentUrl)
	if err != nil {
		log.Debugf("", err)
		log.Warnf("解析torrent失败，结束此流程")
		return nil
	}
	entity.Ep = make([]*models.AnimeEpEntity, 0, len(torrentInfo.Files))

	entity.Torrent = &models.AnimeTorrent{
		Hash: torrentInfo.Hash,
		Url:  torrentInfo.Url,
	}
	title := opts.Title
	for _, t := range torrentInfo.Files {
		// TODO: 筛选文件
		ep := utils.ParseEp(t.Name)
		if ep <= 0 {
			log.Warnf("解析「%s」集数失败，跳过此文件", t.Name)
			continue
		}
		title = t.Name
		entity.Ep = append(entity.Ep, &models.AnimeEpEntity{
			Ep:  ep,
			Src: t.Path(),
		})
	}
	// ------------------- 解析标题获取季度信息 -------------------
	// 优先tmdb解析的season，如果为0则根据配置使用标题解析season
	if entity.Season == 0 {
		var parsed *models.TitleParsed
		// 从后向前解析种子内视频文件
		for i := len(torrentInfo.Files) - 1; i >= 0; i-- {
			parsed = m.parser.Parse(title)
			if parsed.Season > 0 && len(parsed.SeasonRaw) > 0 {
				// 解析成功，跳出循环
				break
			} else {
				parsed = nil
			}
		}
		if parsed == nil {
			entity.Season = m.defaultSeason(0)
		} else {
			d, _ := json.Marshal(parsed)
			fmt.Println(string(d))
			entity.Season = m.defaultSeason(parsed.Season)
		}
		// 没有设置默认季度，且解析失败，结束流程
		if entity.Season <= 0 {
			return nil
		}
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
	return -1
}

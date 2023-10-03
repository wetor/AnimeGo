package parser

import (
	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/xpath"

	"github.com/wetor/AnimeGo/internal/animego/parser/utils"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/torrent"
)

const DefaultSeason = 1

type Manager struct {
	parser  api.ParserPlugin
	mikan   api.AniSource
	bangumi api.AniSource
}

func NewManager(parser api.ParserPlugin, mikan api.AniSource, bangumi api.AniSource) *Manager {
	return &Manager{
		parser:  parser,
		mikan:   mikan,
		bangumi: bangumi,
	}
}

func (m *Manager) Parse(opts *models.ParseOptions) (entity *models.AnimeEntity, err error) {
	if opts.BangumiID > 0 {
		// ------------------- 通过bangumi获取信息（bangumi -> tmdb） -------------------
		entity, err = m.bangumi.Parse(&models.AnimeParseOptions{
			Input:              opts.BangumiID,
			AnimeParseOverride: opts.AnimeParseOverride,
		})
		if err != nil {
			return nil, errors.Wrap(err, "解析anisource失败，结束此流程")
		}
	} else if len(opts.MikanUrl) > 0 {
		// ------------------- 通过mikan获取信息（mikan -> bangumi -> tmdb） -------------------
		entity, err = m.mikan.Parse(&models.AnimeParseOptions{
			Input:              opts.MikanUrl,
			AnimeParseOverride: opts.AnimeParseOverride,
		})
		if err != nil {
			return nil, errors.Wrap(err, "解析anisource失败，结束此流程")
		}
	}

	// ------------------- 获取并解析torrent信息 -------------------
	torrentInfo, err := torrent.LoadUri(opts.TorrentUrl)
	if err != nil {
		return nil, errors.Wrap(err, "解析torrent失败，结束此流程")
	}
	entity.Ep = make([]*models.AnimeEpEntity, 0, len(torrentInfo.Files))
	entity.Flag = models.AnimeFlagNone
	entity.Torrent = &models.AnimeTorrent{
		Hash: torrentInfo.Hash,
	}
	if torrentInfo.Type == torrent.TypeFile {
		entity.Torrent.File = torrentInfo.Url
	} else {
		entity.Torrent.Url = torrentInfo.Url
	}
	for _, t := range torrentInfo.Files {
		// TODO: 筛选文件
		epEntity := &models.AnimeEpEntity{
			Src: xpath.P(t.Path()),
		}
		if isSp, sp := utils.ParseSp(t.Name); isSp {
			epEntity.Type = models.AnimeEpSpecial
			epEntity.Ep = sp
		} else if ep := utils.ParseEp(t.Name); ep > 0 {
			epEntity.Type = models.AnimeEpNormal
			epEntity.Ep = ep
		} else {
			epEntity.Type = models.AnimeEpUnknown
			entity.Flag |= models.AnimeFlagEpParseFailed
			log.Warnf("解析「%s」集数失败，不进行重命名", t.Name)
			// continue
		}
		entity.Ep = append(entity.Ep, epEntity)
	}
	// ------------------- 解析标题获取季度信息 -------------------
	// 优先tmdb解析的season，如果为0则根据配置使用标题解析season
	if entity.Season == 0 {
		var parsed *models.TitleParsed
		// 从后向前解析种子内视频文件，首先解析的是Title
		title := opts.Title
		for i := len(torrentInfo.Files); i >= 0; {
			parsed, err = m.parser.Parse(title)
			if err == nil {
				if parsed.Season > 0 && len(parsed.SeasonRaw) > 0 {
					// 解析成功，跳出循环
					break
				} else {
					parsed = nil
				}
			}
			i--
			if i >= 0 {
				title = torrentInfo.Files[i].Name
			}
		}
		if parsed == nil {
			entity.Season = m.defaultSeason(0)
		} else {
			entity.Season = m.defaultSeason(parsed.Season)
		}
		// 没有设置默认季度，且解析失败，结束流程
		if entity.Season <= 0 {
			err = errors.WithStack(&exceptions.ErrParseFailed{})
			log.DebugErr(err)
			return nil, err
		}
	}
	return entity, nil
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

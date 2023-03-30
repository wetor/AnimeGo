package filter

import (
	"context"

	filterPlugin "github.com/wetor/AnimeGo/internal/animego/filter/plugin"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

type Manager struct {
	filters   []api.FilterPlugin
	anisource api.AniSource
	manager   api.ManagerDownloader
}

// NewManager
//
//	@Description:
//	@param feed api.Feed
//	@param anisource api.AniSource
//	@return *Manager
func NewManager(anisource api.AniSource, manager api.ManagerDownloader) *Manager {
	m := &Manager{
		filters:   make([]api.FilterPlugin, 0),
		anisource: anisource,
		manager:   manager,
	}
	return m
}

func (m *Manager) Add(pluginInfo *models.Plugin) {
	p := filterPlugin.NewFilterPlugin(pluginInfo)
	m.filters = append(m.filters, p)
}

func (m *Manager) Update(ctx context.Context, items []*models.FeedItem) {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	// 筛选
	if len(items) == 0 {
		return
	}

	for _, f := range m.filters {
		items = f.Filter(items)
	}

	for _, item := range items {
		log.Infof("获取「%s」信息开始...", item.Name)
		err := ParseFeedItem(item)
		if err != nil {
			log.Debug(err)
			log.Warnf("解析下载链接失败，目前仅支持torrent文件和magnet链接，结束此流程")
			continue
		}
		// ------------------- 解析文件名获取ep -------------------
		if item.NameParsed == nil || item.NameParsed.Ep == 0 {
			item.NameParsed = public.ParserName(item.Name)
			if item.NameParsed.Ep == 0 {
				log.Warnf("解析集数信息失败，结束此流程")
				continue
			}
		}
		anime := m.anisource.Parse(&models.AnimeParseOptions{
			Url:    item.Url,
			Ep:     item.NameParsed.Ep,
			Season: item.NameParsed.Season,
		})
		if anime != nil {
			anime.DownloadInfo = &models.DownloadInfo{
				Url:  item.Download,
				Hash: item.Hash,
			}

			log.Debugf("发送 %s 下载项:「%s」", item.DownloadType, anime.FullName())
			// 发送需要下载的信息
			m.manager.Download(anime)
		}
		utils.Sleep(DelaySecond, ctx)
	}
}

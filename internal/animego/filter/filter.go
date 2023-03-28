package filter

import (
	"context"

	filterPlugin "github.com/wetor/AnimeGo/internal/animego/filter/plugin"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
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
		anime := m.anisource.Parse(&models.AnimeParseOptions{
			Url:    item.Url,
			Name:   item.Name,
			Parsed: item.NameParsed,
		})
		if anime != nil {
			anime.DownloadInfo = &models.DownloadInfo{
				Url:  item.Download,
				Hash: item.Hash(),
			}

			log.Debugf("发送下载项:「%s」", anime.FullName())
			// 发送需要下载的信息
			m.manager.Download(anime)
		}
		utils.Sleep(DelaySecond, ctx)
	}
}

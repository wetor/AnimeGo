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
	filters []api.FilterPlugin
	manager api.ManagerDownloader
	parser  api.ParserManager
}

// NewManager
//
//	@Description:
//	@param feed api.Feed
//	@return *Manager
func NewManager(manager api.ManagerDownloader, parser api.ParserManager) *Manager {
	m := &Manager{
		filters: make([]api.FilterPlugin, 0),
		manager: manager,
		parser:  parser,
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
		items = f.FilterAll(items)
	}

	for _, item := range items {
		log.Infof("获取「%s」信息开始...", item.Name)
		anime := m.parser.Parse(&models.ParseOptions{
			Title:      item.Name,
			TorrentUrl: item.Download,
			MikanUrl:   item.Url,
		})
		if anime != nil {
			log.Debugf("发送 %s 下载项:「%s」", item.DownloadType, anime.FullName())
			// 发送需要下载的信息
			m.manager.Download(anime)
		}
		utils.Sleep(DelaySecond, ctx)
	}
}

package filter

import (
	"context"

	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

var Set = wire.NewSet(
	NewManager,
	wire.Bind(new(api.FilterManager), new(*Manager)),
)

type Manager struct {
	filters []api.FilterPlugin
	manager api.ManagerDownloader
	parser  api.ParserManager

	*models.FilterOptions
}

// NewManager
//
//	@Description:
//	@param feed api.Feed
//	@return *Manager
func NewManager(opts *models.FilterOptions, manager api.ManagerDownloader, parser api.ParserManager) *Manager {
	m := &Manager{
		filters:       make([]api.FilterPlugin, 0),
		manager:       manager,
		parser:        parser,
		FilterOptions: opts,
	}
	return m
}

func (m *Manager) Add(pluginInfo *models.Plugin) {
	p := NewFilterPlugin(pluginInfo)
	m.filters = append(m.filters, p)
}

func (m *Manager) Update(ctx context.Context, items []*models.FeedItem,
	skipFilter, skipDelay bool) (err error) {
	// 筛选
	if len(items) == 0 {
		return nil
	}
	if !skipFilter {
		for _, f := range m.filters {
			filterResult, err := f.FilterAll(items)
			if err != nil {
				return err
			}
			items = filterResult
		}
	}

	for _, item := range items {
		log.Infof("获取「%s」信息开始...", item.Name)
		anime, err := m.parser.Parse(&models.ParseOptions{
			Title:              item.Name,
			TorrentUrl:         item.TorrentUrl,
			MikanUrl:           item.MikanUrl,
			BangumiID:          item.BangumiID,
			AnimeParseOverride: item.ParseOverride,
		})
		if err != nil {
			log.Warnf("%s", err)
			continue
		}
		log.Debugf("发送下载项:「%s」", anime.FullName())
		// 发送需要下载的信息
		err = m.manager.Download(anime)
		if err != nil {
			if !exceptions.IsExist(err) {
				return err
			}
		}
		if !skipDelay {
			utils.Sleep(m.DelaySecond, ctx)
		}
	}
	return nil
}

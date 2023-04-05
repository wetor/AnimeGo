package parser

import (
	parserPlugin "github.com/wetor/AnimeGo/internal/animego/parser/plugin"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
)

type Manager struct {
	parser api.ParserPlugin
}

func NewManager(plugin *models.Plugin) *Manager {
	single := false
	if plugin.Type == constant.PluginTypeBuiltin {
		single = true
	}
	return &Manager{
		parser: parserPlugin.NewParserPlugin(plugin, single),
	}
}

func (m *Manager) Parse(title, url string) *models.AnimeMultipleEntity {
	return nil
}

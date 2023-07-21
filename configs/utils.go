package configs

import (
	"github.com/jinzhu/copier"
	"github.com/wetor/AnimeGo/internal/models"
)

func ConvertPluginInfo(info []PluginInfo) []models.Plugin {
	plugins := make([]models.Plugin, len(info))
	_ = copier.Copy(&plugins, &info)
	return plugins
}

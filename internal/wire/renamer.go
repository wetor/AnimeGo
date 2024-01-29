//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/animego/renamer"
	"github.com/wetor/AnimeGo/internal/models"
)

func GetRenamePlugin(plugin *models.Plugin) *renamer.Rename {
	wire.Build(
		renamer.PluginSet,
	)
	return nil
}

func GetRenamer(options *renamer.Options, plugin *models.Plugin) *renamer.Manager {
	wire.Build(
		renamer.PluginSet,
		renamer.Set,
	)
	return nil
}

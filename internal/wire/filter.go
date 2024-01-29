//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource/themoviedb"
	"github.com/wetor/AnimeGo/internal/animego/filter"
	"github.com/wetor/AnimeGo/internal/animego/parser"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
)

func GetFilter(opts *models.FilterOptions, manager api.ManagerDownloader,
	parserOpts *models.ParserOptions, plugin *models.Plugin,
	mikanOpts *mikan.Options, bgmOpts *bangumi.Options, tmdbOpts *themoviedb.Options) *filter.Manager {
	wire.Build(
		parser.PluginSet,
		mikan.Set,
		bangumi.Set,
		themoviedb.Set,
		anisource.BangumiSet,
		anisource.MikanSet,
		parser.Set,
		filter.Set,
	)
	return nil
}

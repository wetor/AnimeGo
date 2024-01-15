//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/animego/anisource/bangumi"
	"github.com/wetor/AnimeGo/internal/animego/anisource/mikan"
	"github.com/wetor/AnimeGo/internal/animego/anisource/themoviedb"
)

func GetMikanData(opts *mikan.Options) *mikan.Mikan {
	wire.Build(
		mikan.Set,
	)
	return nil
}

func GetBangumiData(opts *bangumi.Options) *bangumi.Bangumi {
	wire.Build(
		bangumi.Set,
	)
	return nil
}

func GetThemoviedbData(opts *themoviedb.Options) *themoviedb.Themoviedb {
	wire.Build(
		themoviedb.Set,
	)
	return nil
}

func GetMikan(mikanOpts *mikan.Options, bgmOpts *bangumi.Options, tmdbOpts *themoviedb.Options) *anisource.Mikan {
	wire.Build(
		mikan.Set,
		bangumi.Set,
		themoviedb.Set,
		anisource.BangumiSet,
		anisource.MikanSet,
	)
	return nil
}

func GetBangumi(bgmOpts *bangumi.Options, tmdbOpts *themoviedb.Options) *anisource.Bangumi {
	wire.Build(
		bangumi.Set,
		themoviedb.Set,
		anisource.BangumiSet,
	)
	return nil
}

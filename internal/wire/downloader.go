//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/animego/clientnotifier"
	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/api"
)

func GetDownloader(opts *downloader.Options, client api.Client,
	notifyOpts *clientnotifier.Options, db *database.Database, rename api.Renamer) *downloader.Manager {
	wire.Build(
		clientnotifier.Set,
		downloader.Set,
	)
	return nil
}

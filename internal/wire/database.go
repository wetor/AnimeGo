//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/api"
)

func GetDatabase(opts *database.Options, cache api.Cacher) (*database.Database, error) {
	wire.Build(
		database.Set,
	)
	return nil, nil
}

//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/animego/database"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
)

func GetDatabase(opts *models.DatabaseOptions, cache api.Cacher) (*database.Database, error) {
	wire.Build(
		database.Set,
	)
	return nil, nil
}

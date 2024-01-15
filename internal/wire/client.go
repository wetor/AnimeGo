//go:build wireinject

package wire

import (
	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/client"
	"github.com/wetor/AnimeGo/internal/client/qbittorrent"
	"github.com/wetor/AnimeGo/internal/client/transmission"
	"github.com/wetor/AnimeGo/internal/models"
)

func GetClient(name string, opts *models.ClientOptions) api.Client {
	wire.Build(
		client.Set,
	)
	return nil
}

func GetQBittorrent(opts *models.ClientOptions) *qbittorrent.QBittorrent {
	wire.Build(
		qbittorrent.Set,
	)
	return nil
}

func GetTransmission(opts *models.ClientOptions) *transmission.Transmission {
	wire.Build(
		transmission.Set,
	)
	return nil
}

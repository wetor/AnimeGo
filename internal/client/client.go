package client

import (
	"strings"

	"github.com/google/wire"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/client/qbittorrent"
	"github.com/wetor/AnimeGo/internal/client/transmission"
	"github.com/wetor/AnimeGo/internal/models"
)

var Set = wire.NewSet(
	NewClient,
)

func NewClient(name string, opts *models.ClientOptions, cache api.Cacher) api.Client {
	var c api.Client
	switch strings.ToLower(name) {
	case "qbittorrent":
		c = qbittorrent.NewQBittorrent(opts)
	case "transmission":
		c = transmission.NewTransmission(opts, cache)
	}
	return c
}

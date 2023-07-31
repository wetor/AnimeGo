package api

import (
	"context"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/client"
)

type Client interface {
	Connected() bool
	Name() string
	Start(ctx context.Context)
	Config() *client.Config
	List(opt *client.ListOptions) ([]*client.TorrentItem, error)
	Add(opt *client.AddOptions) error
	Delete(opt *client.DeleteOptions) error
}

type ClientNotifier interface {
	// OnDownloadStart will be sent when a download is started.
	OnDownloadStart([]models.ClientEvent)
	// OnDownloadPause will be sent when a download is paused.
	OnDownloadPause([]models.ClientEvent)
	// OnDownloadStop will be sent when a download is stopped by the user.
	OnDownloadStop([]models.ClientEvent)
	// OnDownloadSeeding will be sent when a torrent download is complete but seeding is still going on.
	OnDownloadSeeding([]models.ClientEvent)
	// OnDownloadComplete will be sent when a download is complete. For BitTorrent downloads, this notification is sent when the download is complete and seeding is over.
	OnDownloadComplete([]models.ClientEvent)
	// OnDownloadError will be sent when a download is stopped due to an error.
	OnDownloadError([]models.ClientEvent)
}

package api

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type DownloaderDeleter interface {
	Delete(hash string) error
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

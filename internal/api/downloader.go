package api

import (
	"context"

	"github.com/wetor/AnimeGo/internal/models"
)

type Downloader interface {
	Connected() bool
	Start(ctx context.Context)
	Config() *models.ClientConfig
	List(opt *models.ClientListOptions) []*models.TorrentItem
	Add(opt *models.ClientAddOptions)
	Delete(opt *models.ClientDeleteOptions)
	GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem
}

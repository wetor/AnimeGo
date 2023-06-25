package api

import (
	"context"

	"github.com/wetor/AnimeGo/internal/models"
)

type Downloader interface {
	Connected() bool
	Start(ctx context.Context)
	List(opt *models.ClientListOptions) ([]*models.TorrentItem, error)
	Add(opt *models.ClientAddOptions) error
	Delete(opt *models.ClientDeleteOptions) error
}

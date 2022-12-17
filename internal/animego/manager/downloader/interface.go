package downloader

import (
	"context"
	"github.com/wetor/AnimeGo/internal/models"
)

type Downloader interface {
	Download(anime *models.AnimeEntity)
	Start(ctx context.Context)
}

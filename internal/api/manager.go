package api

import "github.com/wetor/AnimeGo/internal/models"

type ManagerDownloader interface {
	Download(*models.AnimeEntity) error
}

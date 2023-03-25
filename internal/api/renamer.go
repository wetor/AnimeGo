package api

import "github.com/wetor/AnimeGo/internal/models"

type Renamer interface {
	Rename(*models.RenameOptions)
}

type RenamerPlugin interface {
	Rename(anime *models.AnimeEntity, src string) string
}

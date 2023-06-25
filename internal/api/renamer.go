package api

import "github.com/wetor/AnimeGo/internal/models"

type Renamer interface {
	HasRenameTask(name string) bool
	AddRenameTask(*models.RenameOptions) error
}

type RenamerPlugin interface {
	Rename(anime *models.AnimeEntity, index int, src string) (*models.RenameResult, error)
}

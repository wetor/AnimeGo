package api

import "github.com/wetor/AnimeGo/internal/models"

type Renamer interface {
	AddRenameTask(*models.RenameOptions) error
	HasRenameTask(name string) bool
	SetDownloadState(name string, epIndex int, state models.TorrentState) error
	GetEpTaskState(name string, epIndex int) (int, error)
	GetRenameTaskState(name string) (int, error)
}

type RenamerPlugin interface {
	Rename(anime *models.AnimeEntity, index int, src string) (*models.RenameResult, error)
}

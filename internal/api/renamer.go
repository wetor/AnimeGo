package api

import (
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
)

type Renamer interface {
	Init()
	AddRenameTask(*models.RenameOptions) (*models.RenameAllResult, error)
	HasRenameTask(keys []string) bool
	EnableTask(keys []string) error
	DeleteTask(keys []string)
	SetDownloadState(keys []string, state constant.TorrentState) error
	GetEpTaskState(key string) (int, error)
	GetRenameTaskState(keys []string) (int, error)
}

type RenamerPlugin interface {
	Rename(anime *models.AnimeEntity, index int, src string) (*models.RenameResult, error)
}

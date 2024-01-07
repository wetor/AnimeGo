package api

import "github.com/wetor/AnimeGo/internal/models"

type Database interface {
	Scan() error
	IsExist(data any) bool
	Add(data any) error
	Delete(data any) error
	GetAnimeEntity(hash string) (*models.AnimeEntity, error)
	GetAnimeEntityByName(name string) (*models.AnimeEntity, error)
}

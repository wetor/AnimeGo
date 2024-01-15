package api

import (
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
)

type Client interface {
	Connected() bool
	Name() string
	Start()
	Config() *models.Config
	List(opt *models.ListOptions) ([]*models.TorrentItem, error)
	Add(opt *models.AddOptions) error
	Delete(opt *models.DeleteOptions) error
	State(state string) constant.TorrentState
}

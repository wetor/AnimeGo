package api

import (
	"github.com/wetor/AnimeGo/pkg/client"
)

type Client interface {
	Connected() bool
	Name() string
	Start()
	Config() *client.Config
	List(opt *client.ListOptions) ([]*client.TorrentItem, error)
	Add(opt *client.AddOptions) error
	Delete(opt *client.DeleteOptions) error
	State(state string) client.TorrentState
}

package api

import (
	"context"
	"github.com/wetor/AnimeGo/pkg/client"
)

type Client interface {
	Connected() bool
	Name() string
	Start(ctx context.Context)
	Config() *client.Config
	List(opt *client.ListOptions) ([]*client.TorrentItem, error)
	Add(opt *client.AddOptions) error
	Delete(opt *client.DeleteOptions) error
}

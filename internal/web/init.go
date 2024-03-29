package web

import (
	"sync"

	"github.com/wetor/AnimeGo/internal/web/api"
	"github.com/wetor/AnimeGo/internal/web/websocket"
)

var (
	Debug bool
	Host  string
	Port  int
	WG    *sync.WaitGroup

	API *api.Api
	WS  *websocket.WebSocket
)

type Options struct {
	ApiOptions       *api.Options
	WebSocketOptions *websocket.Options
	Host             string
	Port             int
	WG               *sync.WaitGroup
	Notify           chan []byte

	Debug bool
}

func Init(opts *Options) {
	Debug = opts.Debug
	Host = opts.Host
	Port = opts.Port
	WG = opts.WG
	API = api.NewApi(opts.ApiOptions)
	websocket.Init(opts.WebSocketOptions)
	WS = websocket.NewWebSocket()
}

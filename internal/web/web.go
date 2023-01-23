package web

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wetor/AnimeGo/internal/web/api"
)

var (
	Host string
	Port int
	WG   *sync.WaitGroup
)

type Options struct {
	*api.Options
	Host string
	Port int
	WG   *sync.WaitGroup

	Debug bool
}

func Init(opts *Options) {
	if opts.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	Host = opts.Host
	Port = opts.Port
	WG = opts.WG
	api.Init(opts.Options)
}

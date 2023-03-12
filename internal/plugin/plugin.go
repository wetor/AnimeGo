package plugin

import (
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

var (
	Path string
)

type Options struct {
	Path string
}

func Init(opts *Options) {
	gpython.Init()
	lib.Init()
	plugin.Init(&plugin.Options{
		Path: opts.Path,
	})
}

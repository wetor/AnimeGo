package plugin

import (
	"strings"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/plugin/python"
)

type Options struct {
	Type string
	New  bool
}

func GetPlugin(opts *Options) api.Plugin {
	switch strings.ToLower(opts.Type) {
	case "py", "python":
		return &python.Python{}
	default:
		return &python.Python{}
	}
}

package plugin

import (
	"strings"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/internal/plugin/python"
)

var (
	pluginMap = map[string]Plugin{
		javascript.Type: &javascript.JavaScript{},
		python.Type:     &python.Python{},
	}
)

type Plugin interface {
	Type() string
	Execute(opts *models.PluginExecuteOptions, params models.Object) any
	SetSchema(paramsSchema, resultSchema []string)
}

func GetPlugin(t string) Plugin {
	switch strings.ToLower(t) {
	case "js", "javascript":
		return pluginMap[javascript.Type]
	case "py", "python":
		return pluginMap[python.Type]
	default:
		panic("不支持的 plugin.Type: " + t)
	}
}

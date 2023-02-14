package plugin

import (
	"strings"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/plugin/javascript"
	"github.com/wetor/AnimeGo/internal/plugin/python"
)

const (
	Filter   = "filter"
	Schedule = "schedule"
)

var (
	pluginMap = map[string]map[string]api.Plugin{
		Filter: {
			javascript.Type: &javascript.JavaScriptAdapter{},
			python.Type:     &python.Python{},
		},
		Schedule: {
			javascript.Type: &javascript.JavaScriptAdapter{},
			python.Type:     &python.Python{},
		},
	}
)

func GetPlugin(t string, instanceName string) api.Plugin {
	switch strings.ToLower(t) {
	case "js", "javascript":
		return pluginMap[instanceName][javascript.Type]
	case "py", "python":
		return pluginMap[instanceName][python.Type]
	default:
		panic("不支持的 plugin.Type: " + t)
	}
}

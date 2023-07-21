package api

import (
	"github.com/wetor/AnimeGo/pkg/plugin"
)

type Plugin interface {
	Type() string
	PluginLoader
	PluginRunner
	PluginVariable
}

type PluginLoader interface {
	Load(opts *plugin.LoadOptions)
}

type PluginRunner interface {
	Run(function string, args map[string]any) map[string]any
}

type PluginVariable interface {
	Get(varName string) any
	Set(varName string, val any)
}

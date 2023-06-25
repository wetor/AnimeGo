package api

import (
	"github.com/wetor/AnimeGo/pkg/plugin"
)

type Plugin interface {
	IType
	PluginLoader
	PluginRunner
	PluginVariable
}

type PluginLoader interface {
	Load(opts *plugin.LoadOptions) error
}

type PluginRunner interface {
	Run(function string, args map[string]any) (map[string]any, error)
}

type PluginVariable interface {
	Get(varName string) (any, error)
	Set(varName string, val any) error
}

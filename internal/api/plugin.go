package api

import "github.com/wetor/AnimeGo/internal/models"

type Plugin interface {
	Type() string
	PluginLoader
	PluginRunner
	PluginVariable
}

type PluginLoader interface {
	Load(opts *models.PluginLoadOptions)
}

type PluginRunner interface {
	Run(function string, params models.Object) models.Object
}

type PluginVariable interface {
	Get(varName string) interface{}
	Set(varName string, val interface{})
}

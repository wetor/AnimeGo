package api

import "github.com/wetor/AnimeGo/internal/models"

type Plugin interface {
	Type() string
	Execute(opts *models.PluginExecuteOptions, params models.Object) any
	SetSchema(paramsSchema, resultSchema []string)
}

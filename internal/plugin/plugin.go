package plugin

import (
	"github.com/wetor/AnimeGo/internal/models"
)

type Plugin interface {
	Execute(file string, params models.Object) any
	SetSchema(paramsSchema, resultSchema []string)
}

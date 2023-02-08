package python

import (
	"github.com/wetor/AnimeGo/internal/models"
	pluginutils "github.com/wetor/AnimeGo/internal/plugin/utils"
)

type Function struct {
	ParamsSchema    []*pluginutils.Schema
	ResultSchema    []*pluginutils.Schema
	Name            string
	SkipSchemaCheck bool
	Func            func(params models.Object) models.Object
}

func (f *Function) Run(params models.Object) models.Object {
	if !f.SkipSchemaCheck {
		pluginutils.CheckSchema(f.ParamsSchema, params)
	}

	result := f.Func(params)

	if !f.SkipSchemaCheck {
		pluginutils.CheckSchema(f.ResultSchema, result)
	}
	return result
}

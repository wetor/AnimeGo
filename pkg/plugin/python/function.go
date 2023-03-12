package python

import (
	"github.com/wetor/AnimeGo/pkg/plugin"
)

type Function struct {
	ParamsSchema    []*plugin.Schema
	ResultSchema    []*plugin.Schema
	Name            string
	SkipSchemaCheck bool
	Func            func(args map[string]any) map[string]any
}

func (f *Function) Run(args map[string]any) map[string]any {
	if !f.SkipSchemaCheck {
		plugin.CheckSchema(f.ParamsSchema, args)
	}

	result := f.Func(args)

	if !f.SkipSchemaCheck {
		plugin.CheckSchema(f.ResultSchema, result)
	}
	return result
}

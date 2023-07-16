package python

import (
	"github.com/wetor/AnimeGo/pkg/plugin"
)

type Function struct {
	Exist           bool
	ParamsSchema    []*plugin.Schema
	ResultSchema    []*plugin.Schema
	Name            string
	SkipSchemaCheck bool
	DefaultArgs     map[string]any
	Func            func(args map[string]any) (map[string]any, error)
}

func (f *Function) Run(args map[string]any) (map[string]any, error) {
	if !f.SkipSchemaCheck {
		err := plugin.CheckSchema(f.ParamsSchema, args)
		if err != nil {
			return nil, err
		}
	}

	result, err := f.Func(args)
	if err != nil {
		return nil, err
	}
	if !f.SkipSchemaCheck {
		err = plugin.CheckSchema(f.ResultSchema, result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

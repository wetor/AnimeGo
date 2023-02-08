package javascript

import "github.com/wetor/AnimeGo/internal/models"

type JavaScriptAdapter struct {
	js              *JavaScript
	file            string
	skipSchemaCheck bool
}

func (p *JavaScriptAdapter) Load(opts *models.PluginLoadOptions) {
	p.js = &JavaScript{}
	p.js.SetSchema(opts.Functions[0].ParamsSchema, opts.Functions[0].ResultSchema)
	p.file = opts.File
	p.skipSchemaCheck = opts.Functions[0].SkipSchemaCheck
}

func (p *JavaScriptAdapter) Run(function string, params models.Object) (result models.Object) {
	res := p.js.Execute(&models.PluginExecuteOptions{
		File:            p.file,
		SkipSchemaCheck: p.skipSchemaCheck,
	}, params)
	result = res.(models.Object)
	return
}

func (p *JavaScriptAdapter) Type() string {
	return p.js.Type()
}

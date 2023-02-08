package python

import (
	"os"
	"strings"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"

	"github.com/wetor/AnimeGo/internal/models"
	pyutils "github.com/wetor/AnimeGo/internal/plugin/python/utils"
	pluginutils "github.com/wetor/AnimeGo/internal/plugin/utils"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

const Type = "python"

type Python struct {
	functions map[string]*Function
	ctx       py.Context
}

func (p *Python) preExecute(file string) {
	if p.ctx == nil {
		p.ctx = py.NewContext(py.DefaultContextOpts())
		code, err := os.ReadFile(file)
		if err != nil {
			errors.NewAniErrorD(err).TryPanic()
		}
		codeStr := strings.ReplaceAll(string(code), "\r\n", "\n")
		err = os.WriteFile(file, []byte(codeStr), os.ModePerm)
		if err != nil {
			errors.NewAniErrorD(err).TryPanic()
		}
	}
}

func (p *Python) execute(file string) {
	module, err := py.RunFile(p.ctx, file, py.CompileOpts{
		CurDir: "/",
	}, nil)
	if err != nil {
		py.TracebackDump(err)
		errors.NewAniErrorD(err).TryPanic()
	}

	for name, function := range p.functions {
		function.Func = func(params models.Object) models.Object {
			pyObj := pyutils.Value2PyObject(params)
			res, err := module.Call(name, py.Tuple{pyObj}, nil)
			if err != nil {
				py.TracebackDump(err)
			}
			obj, ok := pyutils.PyObject2Value(res).(models.Object)
			if !ok {
				obj = models.Object{
					"result": obj,
				}
			}
			return obj
		}
	}
}

func (p *Python) endExecute() {

}

func (p *Python) Type() string {
	return Type
}

func (p *Python) Load(opts *models.PluginLoadOptions) {
	p.functions = make(map[string]*Function, len(opts.Functions))
	for _, f := range opts.Functions {
		p.functions[f.Name] = &Function{
			ParamsSchema:    pluginutils.ParseSchemas(f.ParamsSchema),
			ResultSchema:    pluginutils.ParseSchemas(f.ResultSchema),
			Name:            f.Name,
			SkipSchemaCheck: f.SkipSchemaCheck,
		}
	}

	try.This(func() {
		p.preExecute(opts.File)

		file := utils.FindScript(opts.File, models.PyExt)
		p.execute(file)

		p.endExecute()

	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本运行时出错", p.Type())
		log.Debugf("", err)
	})
}

func (p *Python) Run(function string, params models.Object) (result models.Object) {
	try.This(func() {
		f := p.functions[function]
		if !f.SkipSchemaCheck {
			pluginutils.CheckSchema(f.ParamsSchema, params)
		}
		result = p.functions[function].Run(params)

		if !f.SkipSchemaCheck {
			pluginutils.CheckSchema(f.ResultSchema, result)
		}
	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本函数 %s 运行时出错", p.Type(), function)
		log.Debugf("", err)
	})
	return result
}

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
	paramsSchema [][]string
	resultSchema [][]string
	ctx          py.Context
	main         func(params models.Object) models.Object // 主函数
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
	p.main = func(params models.Object) models.Object {
		pyObj := pyutils.Value2PyObject(params)
		res, err := module.Call("main", py.Tuple{pyObj}, nil)
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

func (p *Python) endExecute() {

}

func (p *Python) SetSchema(paramsSchema, resultSchema []string) {
	p.paramsSchema = make([][]string, len(paramsSchema))
	for i, param := range paramsSchema {
		p.paramsSchema[i] = strings.Split(param, ":")
	}

	p.resultSchema = make([][]string, len(resultSchema))
	for i, param := range resultSchema {
		p.resultSchema[i] = strings.Split(param, ":")
	}
}

func (p *Python) Type() string {
	return Type
}

func (p *Python) Execute(opts *models.PluginExecuteOptions, params models.Object) (result any) {
	try.This(func() {
		if !opts.SkipCheck {
			pluginutils.CheckParams(p.paramsSchema, params)
		}
		p.preExecute(opts.File)

		file := utils.FindScript(opts.File, models.PyExt)
		p.execute(file)

		p.endExecute()

		result = p.main(params)
		if !opts.SkipCheck {
			pluginutils.CheckResult(p.resultSchema, result)
		}
	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本运行时出错", p.Type())
		log.Debugf("", err)
	})
	return result
}

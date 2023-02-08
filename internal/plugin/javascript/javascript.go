package javascript

import (
	"os"
	"path"
	"strings"

	"github.com/dop251/goja"

	"github.com/wetor/AnimeGo/internal/models"
	pluginutils "github.com/wetor/AnimeGo/internal/plugin/utils"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
)

const Type = "javascript"

type JavaScript struct {
	*goja.Runtime
	main         func(models.Object) models.Object // 主函数
	paramsSchema []*pluginutils.Schema
	resultSchema []*pluginutils.Schema
}

func (p *JavaScript) preExecute() {
	if p.Runtime == nil {
		p.Runtime = goja.New()
		p.registerFunc()
		_, err := p.RunScript(animeGoBaseFilename, animeGoBaseJs)
		errors.NewAniErrorD(err).TryPanic()
	}
}

func (p *JavaScript) execute(file string) {
	raw, err := os.ReadFile(file)
	errors.NewAniErrorD(err).TryPanic()

	currRootPath = path.Dir(file)
	_, currName = path.Split(file)
	currName = strings.TrimSuffix(currName, path.Ext(file))

	_, err = p.RunScript(file, string(raw))
	errors.NewAniErrorD(err).TryPanic()

	err = p.ExportTo(p.Get(funcMain), &p.main)
	errors.NewAniErrorD(err).TryPanic()
}

func (p *JavaScript) endExecute() {
	p.registerVar()
}

func (p *JavaScript) registerFunc() {
	funcMap := p.initFunc()
	for name, method := range funcMap {
		err := p.Set(name, method)
		errors.NewAniErrorD(err).TryPanic()
	}
}

func (p *JavaScript) registerVar() {
	varMap := p.initVar()
	for name, v := range varMap {
		err := p.Set(name, v)
		errors.NewAniErrorD(err).TryPanic()
	}
}

func (p *JavaScript) SetSchema(paramsSchema, resultSchema []string) {
	p.paramsSchema = pluginutils.ParseSchemas(paramsSchema)
	p.resultSchema = pluginutils.ParseSchemas(resultSchema)
}

func (p *JavaScript) Type() string {
	return Type
}

func (p *JavaScript) Execute(opts *models.PluginExecuteOptions, params models.Object) (result any) {
	try.This(func() {
		if !opts.SkipSchemaCheck {
			pluginutils.CheckSchema(p.paramsSchema, params)
		}
		p.preExecute()

		file := utils.FindScript(opts.File, models.JSExt)
		p.execute(file)

		p.endExecute()

		result = p.main(params)
		if !opts.SkipSchemaCheck {
			pluginutils.CheckSchema(p.resultSchema, result)
		}
	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本运行时出错", p.Type())
		log.Debugf("", err)
	})
	return result
}

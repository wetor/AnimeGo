package python

import (
	"os"
	"strings"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	pyutils "github.com/wetor/AnimeGo/internal/plugin/python/utils"
	pluginutils "github.com/wetor/AnimeGo/internal/plugin/utils"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/try"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const Type = "python"

type Python struct {
	functions map[string]*Function
	variables map[string]*Variable
	ctx       py.Context
	module    *py.Module
	name      string
	dir       string
	file      string
}

func (p *Python) preExecute() {
	if p.ctx == nil {
		p.ctx = py.NewContext(py.ContextOpts{
			SysPaths: []string{p.dir},
		})
		code, err := os.ReadFile(p.file)
		if err != nil {
			errors.NewAniErrorD(err).TryPanic()
		}
		codeStr := strings.ReplaceAll(string(code), "\r\n", "\n")
		err = os.WriteFile(p.file, []byte(codeStr), os.ModePerm)
		if err != nil {
			errors.NewAniErrorD(err).TryPanic()
		}
	}
}

func (p *Python) execute() {
	var err error
	p.module, err = py.RunFile(p.ctx, p.file, py.CompileOpts{
		CurDir: "/",
	}, nil)
	if err != nil {
		py.TracebackDump(err)
		errors.NewAniErrorD(err).TryPanic()
	}

}

func (p *Python) endExecute() {
	for name, function := range p.functions {
		function.Func = func(args models.Object) models.Object {
			pyObj := pyutils.Value2PyObject(args)
			res, err := p.module.Call(name, py.Tuple{pyObj}, nil)
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

	for name, variable := range p.variables {
		_, has := p.module.Globals[name]
		if !has && !variable.Nullable {
			log.Warnf("未找到全局变量 %s", name)
			errors.NewAniErrorf("未找到全局变量 %s", name).TryPanic()
		}
	}

	p.module.Globals["__plugin_name__"] = py.String(p.name)
	p.module.Globals["__plugin_dir__"] = py.String(p.dir)
	p.module.Globals["__animego_version__"] = py.String(os.Getenv("ANIMEGO_VERSION"))
	p.module.Globals["_get_config"] = py.MustNewMethod("_get_config", func(self py.Object, args py.Tuple) (py.Object, error) {
		result := models.Object{}
		yamlFile := xpath.Join(p.dir, p.name+".yaml")
		jsonFile := xpath.Join(p.dir, p.name+".json")
		if utils.IsExist(yamlFile) {
			data, err := os.ReadFile(yamlFile)
			if err != nil {
				return nil, err
			}
			err = yaml.Unmarshal(data, &result)
			if err != nil {
				return nil, err
			}
		} else if utils.IsExist(jsonFile) {
			data, err := os.ReadFile(jsonFile)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(data, &result)
			if err != nil {
				return nil, err
			}
		}
		return pyutils.Value2PyObject(result), nil
	}, 0, `_get_config() -> dict`)

}

func (p *Python) Get(name string) any {
	return pyutils.PyObject2Value(p.module.Globals[name])
}

func (p *Python) Set(name string, val any) {
	p.module.Globals[name] = pyutils.Value2PyObject(val)
}

func (p *Python) Type() string {
	return Type
}

func (p *Python) loadPre(file string) {
	if xpath.IsAbs(file) {
		p.file = xpath.Abs(xpath.P(file))
	} else {
		p.file = xpath.Abs(xpath.Join(constant.PluginPath, xpath.P(file)))
	}
	p.file = utils.FindScript(p.file, ".py")
	p.dir, p.name = xpath.Split(p.file)
	p.name = strings.TrimSuffix(p.name, xpath.Ext(p.file))
}

func (p *Python) Load(opts *models.PluginLoadOptions) {
	p.loadPre(opts.File)
	p.functions = make(map[string]*Function, len(opts.Functions))
	for _, f := range opts.Functions {
		p.functions[f.Name] = &Function{
			ParamsSchema:    pluginutils.ParseSchemas(f.ParamsSchema),
			ResultSchema:    pluginutils.ParseSchemas(f.ResultSchema),
			Name:            f.Name,
			SkipSchemaCheck: f.SkipSchemaCheck,
		}
	}
	p.variables = make(map[string]*Variable, len(opts.Variables))
	for _, v := range opts.Variables {
		p.variables[v.Name] = &Variable{
			Name:     v.Name,
			Nullable: v.Nullable,
		}
	}

	try.This(func() {
		p.preExecute()
		p.execute()
		p.endExecute()

	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本运行时出错", p.Type())
		log.Debugf("", err)
	})
}

func (p *Python) Run(function string, args models.Object) (result models.Object) {
	try.This(func() {
		f := p.functions[function]
		if !f.SkipSchemaCheck {
			pluginutils.CheckSchema(f.ParamsSchema, args)
		}
		result = p.functions[function].Run(args)

		if !f.SkipSchemaCheck {
			pluginutils.CheckSchema(f.ResultSchema, result)
		}
	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本函数 %s 运行时出错", p.Type(), function)
		log.Debugf("", err)
	})
	return result
}

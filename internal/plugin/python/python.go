package python

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"gopkg.in/yaml.v3"

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
	variables map[string]*Variable
	ctx       py.Context
	module    *py.Module
	name      string
	path      string
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
	var err error
	file, err = filepath.Abs(file)
	if err != nil {
		errors.NewAniErrorD(err).TryPanic()
	}
	p.path = path.Dir(file)
	_, p.name = path.Split(file)
	p.name = strings.TrimSuffix(p.name, path.Ext(file))

	p.module, err = py.RunFile(p.ctx, file, py.CompileOpts{
		CurDir: "/",
	}, nil)
	if err != nil {
		py.TracebackDump(err)
		errors.NewAniErrorD(err).TryPanic()
	}

}

func (p *Python) endExecute() {

	for name, function := range p.functions {
		function.Func = func(params models.Object) models.Object {
			pyObj := pyutils.Value2PyObject(params)
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
		variable.Getter = func(name string) interface{} {
			return pyutils.PyObject2Value(p.module.Globals[name])
		}
		variable.Setter = func(name string, val interface{}) {
			p.module.Globals[name] = pyutils.Value2PyObject(val)
		}
	}

	p.module.Globals["__plugin_name__"] = py.String(p.name)
	p.module.Globals["__plugin_path__"] = py.String(p.path)
	p.module.Globals["__animego_version__"] = py.String(os.Getenv("ANIMEGO_VERSION"))
	p.module.Globals["_get_config"] = py.MustNewMethod("_get_config", func(self py.Object, args py.Tuple) (py.Object, error) {
		result := models.Object{}
		yamlFile := path.Join(p.path, p.name+".yaml")
		jsonFile := path.Join(p.path, p.name+".json")
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

func (p *Python) Get(name string) interface{} {
	return p.variables[name].Getter(name)
}

func (p *Python) Set(name string, val interface{}) {
	p.variables[name].Setter(name, val)
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
	p.variables = make(map[string]*Variable, len(opts.Variables))
	for _, v := range opts.Variables {
		p.variables[v.Name] = &Variable{
			Name:     v.Name,
			Nullable: v.Nullable,
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

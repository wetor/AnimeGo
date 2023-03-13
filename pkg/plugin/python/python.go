package python

import (
	"os"
	"strings"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/try"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const Type = "python"

type Python struct {
	functions  map[string]*Function
	variables  map[string]*Variable
	globalVars map[string]any
	ctx        py.Context
	module     *py.Module
	name       string
	dir        string
	file       string
	code       *string
	_type      string
}

func NewPython(_type string) *Python {
	return &Python{
		_type: _type,
	}
}

// preExecute
//
//	@Description: 前置执行
//	@receiver p
func (p *Python) preExecute() {
	if p.ctx == nil {
		p.ctx = py.NewContext(py.ContextOpts{
			SysPaths: []string{p.dir},
		})
	}
}

// execute
//
//	@Description: 执行脚本
//	@receiver p
func (p *Python) execute() {
	var err error
	if p.code != nil {
		var code *py.Code
		code, err = py.Compile(*p.code, p.file, py.ExecMode, 0, true)
		if err != nil {
			py.TracebackDump(err)
			errors.NewAniErrorD(err).TryPanic()
		}
		p.module, err = py.RunCode(p.ctx, code, "", nil)
	} else {
		p.module, err = py.RunFile(p.ctx, p.file, py.CompileOpts{
			CurDir: "/",
		}, nil)
	}
	if err != nil {
		py.TracebackDump(err)
		errors.NewAniErrorD(err).TryPanic()
	}

}

// endExecute
//
//	@Description: 后置执行，写入变量，获取方法
//	@receiver p
func (p *Python) endExecute() {
	for name, function := range p.functions {
		function.Func = func(args map[string]any) map[string]any {
			pyObj := plugin.Value2PyObject(args)
			res, err := p.module.Call(name, py.Tuple{pyObj}, nil)
			if err != nil {
				py.TracebackDump(err)
			}
			obj, ok := plugin.PyObject2Value(res).(map[string]any)
			if !ok {
				obj = map[string]any{
					"result": obj,
				}
			}
			return obj
		}
	}
	for name, val := range p.globalVars {
		p.Set(name, val)
	}
	for name, variable := range p.variables {
		_, has := p.module.Globals[name]
		if !has && !variable.Nullable {
			log.Warnf("未找到全局变量 %s", name)
			errors.NewAniErrorf("未找到全局变量 %s", name).TryPanic()
		}
	}
	p.Set("__plugin_name__", p.name)
	p.Set("__plugin_dir__", p.dir)
	p.Set("__animego_version__", os.Getenv("ANIMEGO_VERSION"))
	p.Set("_get_config", py.MustNewMethod("_get_config", func(self py.Object, args py.Tuple) (py.Object, error) {
		result := map[string]any{}
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
		return plugin.Value2PyObject(result), nil
	}, 0, `_get_config() -> dict`))

}

// Get
//
//	@Description: 获取变量
//	@receiver p
//	@param name
//	@return any
func (p *Python) Get(name string) any {
	return plugin.PyObject2Value(p.module.Globals[name])
}

// Set
//
//	@Description: 设置变量
//	@receiver p
//	@param name
//	@param val
func (p *Python) Set(name string, val any) {
	p.module.Globals[name] = plugin.Value2PyObject(val)
}

// Type
//
//	@Description: 脚本类型
//	@receiver p
//	@return string
func (p *Python) Type() string {
	return p._type
}

// loadPre
//
//	@Description: 前置加载，脚本路径转为绝对路径
//	@receiver p
//	@param file
func (p *Python) loadPre(file string) {
	if xpath.IsAbs(file) {
		p.file = xpath.Abs(xpath.P(file))
	} else {
		p.file = xpath.Abs(xpath.Join(plugin.Path, xpath.P(file)))
	}
	p.file = utils.FindScript(p.file, ".py")
	p.dir, p.name = xpath.Split(p.file)
	p.name = strings.TrimSuffix(p.name, xpath.Ext(p.file))
}

// Load
//
//	@Description: 加载脚本
//	@receiver p
//	@param opts
func (p *Python) Load(opts *plugin.LoadOptions) {
	p.globalVars = opts.GlobalVars
	if opts.Code == nil {
		p.loadPre(opts.File)
	} else {
		p.code = opts.Code
	}
	p.functions = make(map[string]*Function, len(opts.FuncSchema))
	for _, f := range opts.FuncSchema {
		p.functions[f.Name] = &Function{
			ParamsSchema:    plugin.ParseSchemas(f.ParamsSchema),
			ResultSchema:    plugin.ParseSchemas(f.ResultSchema),
			Name:            f.Name,
			SkipSchemaCheck: f.SkipSchemaCheck,
			DefaultArgs:     f.DefaultArgs,
		}
	}
	p.variables = make(map[string]*Variable, len(opts.VarSchema))
	for _, v := range opts.VarSchema {
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

// Run
//
//	@Description: 执行脚本函数
//	@receiver p
//	@param function 函数名
//	@param args 参数列表
//	@return result
func (p *Python) Run(function string, args map[string]any) (result map[string]any) {
	try.This(func() {
		f := p.functions[function]
		for k, v := range f.DefaultArgs {
			if _, ok := args[k]; !ok {
				args[k] = v
			}
		}
		result = p.functions[function].Run(args)
	}).Catch(func(err try.E) {
		log.Warnf("%s 脚本函数 %s 运行时出错", p.Type(), function)
		log.Debugf("", err)
	})
	return result
}

package python

import (
	"os"
	"path"
	"strings"

	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

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
func (p *Python) execute() (err error) {
	if p.code != nil {
		var code *py.Code
		code, err = py.Compile(*p.code, p.file, py.ExecMode, 0, true)
		if err != nil {
			py.TracebackDump(err)
			log.DebugErr(err)
			log.Warnf("%s: %s", "编译失败", p.file)
			return errors.Wrap(err, "编译失败")
		}
		p.module, err = py.RunCode(p.ctx, code, "", nil)
	} else {
		p.module, err = py.RunFile(p.ctx, p.file, py.CompileOpts{
			CurDir: "/",
		}, nil)
	}
	if err != nil {
		py.TracebackDump(err)
		log.DebugErr(err)
		log.Warnf("%s: %s", "执行失败", p.file)
		return errors.Wrap(err, "执行失败")
	}
	return nil
}

// endExecute
//
//	@Description: 后置执行，写入变量，获取方法
//	@receiver p
func (p *Python) endExecute() (err error) {
	for name, function := range p.functions {
		f, ok := p.module.Globals[name]
		if !ok {
			continue
		}
		caller, ok := f.(py.I__call__)
		if !ok {
			continue
		}
		function.Exist = true
		function.Func = func(args map[string]any) (map[string]any, error) {
			pyObj, err := ToObject(args)
			if err != nil {
				return nil, err
			}
			res, err := caller.M__call__(py.Tuple{pyObj}, nil)
			if err != nil {
				py.TracebackDump(err)
				log.DebugErr(err)
				return nil, errors.Wrapf(err, "函数调用失败: %s", name)
			}
			val, err := ToValue(res)
			if err != nil {
				return nil, err
			}
			obj, ok := val.(map[string]any)
			if !ok {
				obj = map[string]any{
					"result": obj,
				}
			}
			return obj, nil
		}
	}
	for name, val := range p.globalVars {
		err = p.Set(name, val)
		if err != nil {
			return err
		}
	}
	for name, variable := range p.variables {
		_, has := p.module.Globals[name]
		if !has && !variable.Nullable {
			err = errors.WithStack(exceptions.ErrPlugin{Type: p.Type(), File: p.file, Message: "未找到全局变量: " + name})
			log.DebugErr(err)
			return err
		}
	}
	err = p.Set("__debug__", plugin.Debug)
	if err != nil {
		return err
	}
	err = p.Set("__plugin_name__", p.name)
	if err != nil {
		return err
	}
	err = p.Set("__plugin_dir__", p.dir)
	if err != nil {
		return err
	}
	err = p.Set("_get_config", py.MustNewMethod("_get_config", func(self py.Object, args py.Tuple) (py.Object, error) {
		result := map[string]any{}
		yamlFile := path.Join(p.dir, p.name+".yaml")
		jsonFile := path.Join(p.dir, p.name+".json")
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
		return ToObject(result)
	}, 0, `_get_config() -> dict`))
	return err
}

// Get
//
//	@Description: 获取变量
//	@receiver p
//	@param name
//	@return any
func (p *Python) Get(name string) (any, error) {
	return ToValue(p.module.Globals[name])
}

// Set
//
//	@Description: 设置变量
//	@receiver p
//	@param name
//	@param val
func (p *Python) Set(name string, val any) error {
	if m, ok := val.(*py.Method); ok {
		p.module.Globals[name] = m
	} else {
		obj, err := ToObject(val)
		if err != nil {
			return err
		}
		p.module.Globals[name] = obj
	}
	return nil
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
func (p *Python) loadPre(file string) (err error) {
	if xpath.IsAbs(file) {
		p.file = xpath.Abs(file)
	} else {
		p.file = xpath.Abs(path.Join(plugin.Path, file))
	}
	p.file, err = utils.FindScript(p.file, ".py")
	if err != nil {
		return errors.Wrap(err, "加载插件失败")
	}
	p.dir, p.name = path.Split(p.file)
	p.name = strings.TrimSuffix(p.name, path.Ext(p.file))
	return nil
}

// Load
//
//	@Description: 加载脚本
//	@receiver p
//	@param opts
func (p *Python) Load(opts *plugin.LoadOptions) (err error) {
	p.globalVars = opts.GlobalVars
	if opts.Code == nil {
		err = p.loadPre(opts.File)
		if err != nil {
			return err
		}
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
	p.preExecute()
	err = p.execute()
	if err != nil {
		return err
	}
	err = p.endExecute()
	if err != nil {
		return err
	}
	return nil
}

// Run
//
//	@Description: 执行脚本函数
//	@receiver p
//	@param function 函数名
//	@param args 参数列表
//	@return result
func (p *Python) Run(function string, args map[string]any) (result map[string]any, err error) {
	f := p.functions[function]
	if !f.Exist {
		log.Warnf("%s 脚本函数 %s 不存在，跳过", p.Type(), function)
		return
	}
	for k, v := range f.DefaultArgs {
		if _, ok := args[k]; !ok {
			args[k] = v
		}
	}
	result, err = p.functions[function].Run(args)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.Wrapf(err, "%s 脚本函数 %s 运行时出错", p.Type(), function)
	}
	return result, nil
}

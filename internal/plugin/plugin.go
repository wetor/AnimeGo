package plugin

import (
	"os"
	"strings"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python/lib"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

type Options struct {
	Path  string
	Debug bool
}

func Init(opts *Options) {
	gpython.Init()
	lib.Init()
	plugin.Init(&plugin.Options{
		Path:  opts.Path,
		Debug: opts.Debug,
	})
}

type LoadPluginOptions struct {
	*models.Plugin
	// EntryFunc
	//  @Description: 入口函数名
	EntryFunc string
	// FuncSchema
	//  @Description: 函数约束
	FuncSchema []*plugin.FuncSchemaOptions
	// VarSchema
	//  @Description: 全局变量约束
	VarSchema []*plugin.VarSchemaOptions
}

func (o *LoadPluginOptions) Default() {
	if o.Args == nil {
		o.Args = make(map[string]any)
	}
	if o.Vars == nil {
		o.Vars = make(map[string]any)
	}
}

func LoadPlugin(opts *LoadPluginOptions) (p api.Plugin) {
	opts.Default()
	var code *string = nil
	pluginType := strings.ToLower(opts.Type)
	switch pluginType {
	case constant.PluginTypePython, "py":
		p = python.NewPython(pluginType)
	case constant.PluginTypeBuiltin:
		p = python.NewPython(pluginType)
		code = assets.GetBuiltinPlugin(opts.File)
	default:
		log.Warnf("不支持的插件类型 %s", pluginType)
		errors.NewAniErrorf("不支持的插件类型 %s", pluginType).TryPanic()
	}
	for _, f := range opts.FuncSchema {
		if f.Name == opts.EntryFunc {
			f.DefaultArgs = opts.Args
			break
		}
	}
	p.Load(&plugin.LoadOptions{
		File:       opts.File,
		Code:       code,
		GlobalVars: opts.Vars,
		FuncSchema: opts.FuncSchema,
		VarSchema:  opts.VarSchema,
	})
	p.Set("__animego_version__", os.Getenv("ANIMEGO_VERSION"))
	return
}

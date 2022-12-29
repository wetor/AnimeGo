package javascript

import (
	"github.com/dop251/goja"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"os"
	"path"
	"strings"
)

type JavaScript struct {
	*goja.Runtime
	main         func(plugin.Object) plugin.Object // 主函数
	paramsSchema []string
	resultSchema []string
}

// preExecute
//  @Description: 前置处理：初始化js虚拟机，注册函数，执行基础脚本
//  @receiver *JavaScript
//
func (js *JavaScript) preExecute() {
	if js.Runtime == nil {
		js.Runtime = goja.New()
		js.registerFunc()
		_, err := js.RunScript(animeGoBaseFilename, animeGoBaseJs)
		errors.NewAniErrorD(err).TryPanic()
	}
}

// execute
//  @Description: 执行脚本
//  @receiver *JavaScript
//  @param file string
//  @return result any
//
func (js *JavaScript) execute(file string) {
	raw, err := os.ReadFile(file)
	errors.NewAniErrorD(err).TryPanic()

	currRootPath = path.Dir(file)
	_, currName = path.Split(file)
	currName = strings.TrimSuffix(currName, path.Ext(file))

	_, err = js.RunScript(file, string(raw))
	errors.NewAniErrorD(err).TryPanic()

	err = js.ExportTo(js.Get(funcMain), &js.main)
	errors.NewAniErrorD(err).TryPanic()
}

// endExecute
//  @Description: 后置处理
//  @receiver *JavaScript
//
func (js *JavaScript) endExecute() {
	js.registerVar()
}

// registerFunc
//  @Description: 注册函数
//  @receiver *JavaScript
//
func (js *JavaScript) registerFunc() {
	funcMap := js.initFunc()
	for name, method := range funcMap {
		err := js.Set(name, method)
		errors.NewAniErrorD(err).TryPanic()
	}
}

// registerVar
//  @Description: 注册全局变量
//  @receiver *JavaScript
//
func (js *JavaScript) registerVar() {
	varMap := js.initVar()
	for name, v := range varMap {
		err := js.Set(name, v)
		errors.NewAniErrorD(err).TryPanic()
	}
}

func (js *JavaScript) checkParams(params plugin.Object) {
	for _, field := range js.paramsSchema {
		_, has := params[field]
		if !has {
			errors.NewAniError("参数缺少: " + field).TryPanic()
		}
	}
}

func (js *JavaScript) checkResult(result any) {
	resultMap, ok := result.(plugin.Object)
	if !ok {
		errors.NewAniError("返回类型错误").TryPanic()
	}
	for _, field := range js.resultSchema {
		_, has := resultMap[field]
		if !has {
			errors.NewAniError("返回值缺少: " + field).TryPanic()
		}
	}
}

func (js *JavaScript) SetSchema(paramsSchema, resultSchema []string) {
	js.paramsSchema = paramsSchema
	js.resultSchema = resultSchema
}

func (js *JavaScript) Execute(file string, params plugin.Object) (result any) {
	js.checkParams(params)

	js.preExecute()

	file = FindScript(file)
	js.execute(file)

	js.endExecute()

	result = js.main(params)
	js.checkResult(result)
	return result
}

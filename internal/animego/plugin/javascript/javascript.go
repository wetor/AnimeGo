package javascript

import (
	"AnimeGo/pkg/errors"
	"github.com/dop251/goja"
	"os"
)

type JavaScript struct {
	*goja.Runtime
	paramsSchema []string
	resultSchema []string
}

// preExecute
//  @Description: 前置处理：初始化js虚拟机，注册函数，执行基础脚本
//  @receiver *JavaScript
//  @return error
//
func (js *JavaScript) preExecute() error {
	if js.Runtime == nil {
		js.Runtime = goja.New()
		err := js.registerFunc()
		if err != nil {
			return err
		}
		_, err = js.RunScript(animeGoBaseFilename, animeGoBaseJs)
		if err != nil {
			return err
		}
	}
	return nil
}

// execute
//  @Description: 执行脚本
//  @receiver *JavaScript
//  @param file string
//  @param params map[string]interface{}
//  @return result interface{}
//  @return err error
//
func (js *JavaScript) execute(file string, params Object) (result interface{}, err error) {
	raw, err := os.ReadFile(file)
	if err != nil {
		return
	}

	_, err = js.RunScript(file, string(raw))
	if err != nil {
		return
	}

	var main func(Object) Object
	err = js.ExportTo(js.Get(funcMain), &main)
	if err != nil {
		return
	}
	result = main(params)
	return
}

// endExecute
//  @Description: 后置处理
//  @receiver *JavaScript
//  @return error
//
func (js *JavaScript) endExecute() error {
	return nil
}

// registerFunc
//  @Description: 注册函数
//  @receiver *JavaScript
//  @return error
//
func (js *JavaScript) registerFunc() error {
	funcMap := js.initFunc()
	for name, method := range funcMap {
		err := js.Set(name, method)
		if err != nil {
			return err
		}
	}
	return nil
}

func (js *JavaScript) checkParams(params Object) error {
	for _, field := range js.paramsSchema {
		_, has := params[field]
		if !has {
			return errors.NewAniError("参数缺少: " + field)
		}
	}
	return nil
}

func (js *JavaScript) checkResult(result interface{}) error {
	resultMap, ok := result.(Object)
	if !ok {
		return errors.NewAniError("返回类型错误")
	}
	for _, field := range js.resultSchema {
		_, has := resultMap[field]
		if !has {
			return errors.NewAniError("返回值缺少: " + field)
		}
	}
	return nil
}

func (js *JavaScript) SetSchema(paramsSchema, resultSchema []string) {
	js.paramsSchema = paramsSchema
	js.resultSchema = resultSchema
}

func (js *JavaScript) Execute(file string, params Object) (result interface{}, err error) {
	err = js.checkParams(params)
	if err != nil {
		return
	}
	err = js.preExecute()
	if err != nil {
		return
	}
	result, err = js.execute(file, params)
	if err != nil {
		return
	}
	err = js.endExecute()
	if err != nil {
		return
	}
	err = js.checkResult(result)
	if err != nil {
		return
	}
	return
}

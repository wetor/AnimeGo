package javascript

import (
	"AnimeGo/pkg/errors"
	"fmt"
	"github.com/dop251/goja"
	"os"
	"path"
	"strings"
)

type JavaScript struct {
	*goja.Runtime
	main         func(Object) Object // 主函数
	rootPath     string              // 脚本所在路径
	name         string              // 脚本名
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
			return errors.NewAniErrorD(err)
		}
		_, err = js.RunScript(animeGoBaseFilename, animeGoBaseJs)
		if err != nil {
			return errors.NewAniErrorD(err)
		}
	}
	return nil
}

// execute
//  @Description: 执行脚本
//  @receiver *JavaScript
//  @param file string
//  @return result interface{}
//  @return err error
//
func (js *JavaScript) execute(file string) error {
	raw, err := os.ReadFile(file)
	if err != nil {
		return errors.NewAniErrorD(err)
	}
	js.rootPath = path.Dir(file)
	_, js.name = path.Split(file)
	js.name = strings.TrimSuffix(js.name, path.Ext(file))

	fmt.Println(js.rootPath, js.name)

	_, err = js.RunScript(file, string(raw))
	if err != nil {
		return errors.NewAniErrorD(err)
	}

	err = js.ExportTo(js.Get(funcMain), &js.main)
	if err != nil {
		return errors.NewAniErrorD(err)
	}
	return nil
}

// endExecute
//  @Description: 后置处理
//  @receiver *JavaScript
//  @return error
//
func (js *JavaScript) endExecute() error {
	err := js.registerVar()
	if err != nil {
		return errors.NewAniErrorD(err)
	}
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

// registerVar
//  @Description: 注册全局变量
//  @receiver *JavaScript
//  @return error
//
func (js *JavaScript) registerVar() error {
	varMap := js.initVar()
	for name, v := range varMap {
		err := js.Set(name, v)
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
	err = js.execute(file)
	if err != nil {
		return
	}
	err = js.endExecute()
	if err != nil {
		return
	}

	result = js.main(params)
	err = js.checkResult(result)
	if err != nil {
		return
	}
	return
}

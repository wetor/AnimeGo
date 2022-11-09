package javascript

import (
	"github.com/dop251/goja"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"os"
	"path"
	"strings"
)

type JavaScript struct {
	*goja.Runtime
	main         func(Object) Object // 主函数
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
	currRootPath = path.Dir(file)
	_, currName = path.Split(file)
	currName = strings.TrimSuffix(currName, path.Ext(file))

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

func (js *JavaScript) findScript(file string) (string, error) {
	if utils.IsDir(file) {
		// 文件夹，在文件夹中寻找 main.js 和 plugin.js
		if utils.IsExist(path.Join(file, "main.js")) {
			return path.Join(file, "main.js"), nil
		} else if utils.IsExist(path.Join(file, "plugin.js")) {
			return path.Join(file, "plugin.js"), nil
		} else {
			return "", errors.NewAniError("插件文件夹中找不到 'main.js' 或 'plugin.js'")
		}
	} else if !utils.IsExist(file) {
		// 文件不存在，尝试增加 .js 扩展名
		if utils.IsExist(file + ".js") {
			return file + ".js", nil
		} else {
			return "", errors.NewAniError("插件文件不存在")
		}
	} else if path.Ext(file) != ".js" {
		// 文件存在，判断扩展名是否为 .js
		return "", errors.NewAniError("插件文件格式错误")
	}
	return file, nil
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
	file, err = js.findScript(file)
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

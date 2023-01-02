package javascript

import (
	"github.com/dop251/goja"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"os"
	"path"
	"strings"
)

type JavaScript struct {
	*goja.Runtime
	main         func(models.Object) models.Object // 主函数
	paramsSchema []string
	resultSchema []string
}

func (js *JavaScript) preExecute() {
	if js.Runtime == nil {
		js.Runtime = goja.New()
		js.registerFunc()
		_, err := js.RunScript(animeGoBaseFilename, animeGoBaseJs)
		errors.NewAniErrorD(err).TryPanic()
	}
}

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

func (js *JavaScript) endExecute() {
	js.registerVar()
}

func (js *JavaScript) registerFunc() {
	funcMap := js.initFunc()
	for name, method := range funcMap {
		err := js.Set(name, method)
		errors.NewAniErrorD(err).TryPanic()
	}
}

func (js *JavaScript) registerVar() {
	varMap := js.initVar()
	for name, v := range varMap {
		err := js.Set(name, v)
		errors.NewAniErrorD(err).TryPanic()
	}
}

func (js *JavaScript) checkParams(params models.Object) {
	for _, field := range js.paramsSchema {
		_, has := params[field]
		if !has {
			errors.NewAniError("参数缺少: " + field).TryPanic()
		}
	}
}

func (js *JavaScript) checkResult(result any) {
	resultMap, ok := result.(models.Object)
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

func (js *JavaScript) Execute(file string, params models.Object) (result any) {
	func() {
		defer errors.HandleError(func(err error) {
			zap.S().Error(err)
		})
		js.checkParams(params)

		js.preExecute()

		file = utils.FindScript(file, models.JSExt)
		js.execute(file)

		js.endExecute()

		result = js.main(params)
		js.checkResult(result)
	}()
	return result
}

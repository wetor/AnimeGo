package javascript

import (
	"os"
	"path"
	"strings"

	"github.com/dop251/goja"
	"go.uber.org/zap"

	"github.com/wetor/AnimeGo/internal/models"
	pluginutils "github.com/wetor/AnimeGo/internal/plugin/utils"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
)

const Type = "javascript"

type JavaScript struct {
	*goja.Runtime
	main         func(models.Object) models.Object // 主函数
	paramsSchema [][]string
	resultSchema [][]string
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

func (js *JavaScript) SetSchema(paramsSchema, resultSchema []string) {
	js.paramsSchema = make([][]string, len(paramsSchema))
	for i, param := range paramsSchema {
		js.paramsSchema[i] = strings.Split(param, ":")
	}

	js.resultSchema = make([][]string, len(resultSchema))
	for i, param := range resultSchema {
		js.resultSchema[i] = strings.Split(param, ":")
	}
}

func (js *JavaScript) Type() string {
	return Type
}

func (js *JavaScript) Execute(file string, params models.Object) (result any) {
	func() {
		defer errors.HandleError(func(err error) {
			zap.S().Error(err)
		})
		pluginutils.CheckParams(js.paramsSchema, params)

		js.preExecute()

		file = utils.FindScript(file, models.JSExt)
		js.execute(file)

		js.endExecute()

		result = js.main(params)
		pluginutils.CheckResult(js.resultSchema, result)
	}()
	return result
}

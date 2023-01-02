package python

import (
	gpy "github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
)

type Python struct {
	paramsSchema []string
	resultSchema []string
	ctx          gpy.Context
	main         func(params models.Object) models.Object // 主函数
}

func (py *Python) preExecute() {
	if py.ctx == nil {
		py.ctx = gpy.NewContext(gpy.DefaultContextOpts())
	}
}

func (py *Python) execute(file string) {
	module, err := gpy.RunFile(py.ctx, file, gpy.CompileOpts{
		CurDir: "/",
	}, nil)
	if err != nil {
		gpy.TracebackDump(err)
		errors.NewAniErrorD(err).TryPanic()
	}
	py.main = func(params models.Object) models.Object {
		pyObj := Value2PyObject(params)
		res, err := module.Call("main", gpy.Tuple{pyObj}, nil)
		if err != nil {
			gpy.TracebackDump(err)
		}
		obj, ok := PyObject2Value(res).(models.Object)
		if !ok {
			obj = models.Object{
				"result": obj,
			}
		}
		return obj
	}
}

func (py *Python) endExecute() {

}

func (py *Python) checkParams(params models.Object) {
	for _, field := range py.paramsSchema {
		_, has := params[field]
		if !has {
			errors.NewAniError("参数缺少: " + field).TryPanic()
		}
	}
}

func (py *Python) checkResult(result any) {
	resultMap, ok := result.(models.Object)
	if !ok {
		errors.NewAniError("返回类型错误").TryPanic()
	}
	for _, field := range py.resultSchema {
		_, has := resultMap[field]
		if !has {
			errors.NewAniError("返回值缺少: " + field).TryPanic()
		}
	}
}

func (py *Python) SetSchema(paramsSchema, resultSchema []string) {
	py.paramsSchema = paramsSchema
	py.resultSchema = resultSchema
}

func (py *Python) Execute(file string, params models.Object) (result any) {
	func() {
		defer errors.HandleError(func(err error) {
			zap.S().Error(err)
		})
		py.checkParams(params)

		py.preExecute()

		file = utils.FindScript(file, models.PyExt)
		py.execute(file)

		py.endExecute()

		result = py.main(params)
		py.checkResult(result)
	}()
	return result
}

package python

import (
	"fmt"
	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
	"github.com/wetor/AnimeGo/internal/plugin"
)

type Python struct {
	paramsSchema []string
	resultSchema []string
}

func (py *Python) SetSchema(paramsSchema, resultSchema []string) {
	py.paramsSchema = paramsSchema
	py.resultSchema = resultSchema
}

func (py *Python) Execute(file string, params plugin.Object) any {

	return nil
}

func RunWithFile(pyFile string, test bool) error {

	// See type Context interface and related docs
	ctx := py.NewContext(py.DefaultContextOpts())

	// This drives modules being able to perform cleanup and release resources
	defer ctx.Close()
	m, err := py.RunFile(ctx, pyFile, py.CompileOpts{}, nil)
	if err != nil {
		py.TracebackDump(err)
	}

	var res py.Object
	if test {
		res, err = m.Call("__test__", nil, nil)
		if err != nil {
			py.TracebackDump(err)
		}
	} else {
		res, err = m.Call("main", py.Tuple{}, nil)
		if err != nil {
			py.TracebackDump(err)
		}
	}

	if res != nil {
		fmt.Println(res)
	}

	return err
}

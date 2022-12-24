package python

import (
	"github.com/go-python/gpython/py"
	_ "github.com/go-python/gpython/stdlib"
)

type Python struct {
}

func (py *Python) Execute(file string, params map[string]any) any {

	return nil
}

func RunWithFile(pyFile string) error {

	// See type Context interface and related docs
	ctx := py.NewContext(py.DefaultContextOpts())

	// This drives modules being able to perform cleanup and release resources
	defer ctx.Close()
	_, err := py.RunFile(ctx, pyFile, py.CompileOpts{}, nil)

	if err != nil {
		py.TracebackDump(err)
	}

	return err
}

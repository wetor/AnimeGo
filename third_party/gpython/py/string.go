package py

import (
	"github.com/go-python/gpython/py"
	"strings"
)

func init() {
	py.StringType.Dict["strip"] = py.MustNewMethod("strip", func(self py.Object) (py.Object, error) {
		return py.String(strings.TrimSpace(string(self.(py.String)))), nil
	}, 0, `strip()`)

	py.StringType.Dict["join"] = py.MustNewMethod("join", func(self py.Object, args py.Object) (py.Object, error) {
		argList := args.(*py.List)
		list := make([]string, argList.Len())
		for i, item := range argList.Items {
			list[i] = string(item.(py.String))
		}
		return py.String(strings.Join(list, string(self.(py.String)))), nil
	}, 0, `join(list)`)

}

package builtin

import (
	"github.com/go-python/gpython/py"
)

func Init() {
	methods := []*py.Method{
		py.MustNewMethod("map", builtin_map, 0, "map"),
		py.MustNewMethod("filter", builtin_filter, 0, "filter"),
	}

	builtins := py.GetModuleImpl("builtins")
	builtins.Methods = append(builtins.Methods, methods...)
}

func _zip(args py.Tuple) (*py.Zip, error) {
	newObj, err := py.ZipTypeNew(nil, args, nil)
	if err != nil {
		return nil, err
	}
	return newObj.(*py.Zip), nil
}

func builtin_map(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	method := args[0].(*py.Function)

	zip, err := _zip(args[1:])
	if err != nil {
		return nil, err
	}
	result := py.NewList()
	for k, _ := zip.M__next__(); k != nil; k, _ = zip.M__next__() {
		tuple := k.(py.Tuple)
		obj, err := method.M__call__(tuple, nil)
		if err != nil {
			return nil, err
		}
		result.Append(obj)
	}
	return py.NewIterator(result.Items), nil
}

func builtin_filter(self py.Object, args py.Tuple, kwargs py.StringDict) (py.Object, error) {
	method := args[0].(*py.Function)

	zip, err := _zip(args[1:])
	if err != nil {
		return nil, err
	}

	result := py.NewList()
	for k, _ := zip.M__next__(); k != nil; k, _ = zip.M__next__() {
		tuple := k.(py.Tuple)
		obj, err := method.M__call__(tuple, nil)
		if err != nil {
			return nil, err
		}
		if obj.(py.Bool) == true {
			if len(tuple) == 1 {
				result.Append(tuple[0])
			} else if len(tuple) > 1 {
				result.Append(tuple)
			}
		}
	}
	return py.NewIterator(result.Items), nil
}

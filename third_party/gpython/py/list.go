package py

import "github.com/go-python/gpython/py"

func init() {
	py.ListType.Dict["remove"] = py.MustNewMethod("remove", func(self py.Object, args py.Object) (py.Object, error) {
		listSelf := self.(*py.List)
		for i, obj := range listSelf.Items {
			if args == obj {
				listSelf.Items = append(listSelf.Items[:i], listSelf.Items[i+1:]...)
				break
			}
		}
		return py.NoneType{}, nil
	}, 0, "remove(item)")

}

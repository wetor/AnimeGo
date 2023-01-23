package python

import (
	"reflect"

	gpy "github.com/go-python/gpython/py"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
)

func Value2PyObject(object any) gpy.Object {
	var pyVal gpy.Object
	switch val := object.(type) {
	case nil:
		pyVal = gpy.None
	case string:
		pyVal = gpy.String(val)
	case int:
		pyVal = gpy.Int(val)
	case int64:
		pyVal = gpy.Int(val)
	case float32:
		pyVal = gpy.Float(val)
	case float64:
		pyVal = gpy.Float(val)
	case bool:
		pyVal = gpy.NewBool(val)
	case models.Object:
		pyValDict := gpy.NewStringDictSized(len(val))
		for key, value := range val {
			pyValDict[key] = Value2PyObject(value)
		}
		pyVal = pyValDict
	default:
		refVal := reflect.ValueOf(object)
		switch refVal.Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			l := refVal.Len()
			pyValList := gpy.NewListWithCapacity(l)
			for i := 0; i < l; i++ {
				pyValList.Append(Value2PyObject(refVal.Index(i).Interface()))
			}
			pyVal = pyValList
		default:
			errors.NewAniError("不支持的类型").TryPanic()
		}
	}
	return pyVal
}

func PyObject2Value(object gpy.Object) any {
	var objVal any
	switch val := object.(type) {
	case gpy.String:
		objVal = string(val)
	case gpy.Int:
		objVal = int64(val)
	case gpy.Bool:
		objVal = bool(val)
	case gpy.Float:
		objVal = float64(val)
	case gpy.StringDict:
		objValDict := make(models.Object, len(val))
		for key, value := range val {
			objValDict[key] = PyObject2Value(value)
		}
		objVal = objValDict
	case *gpy.List:
		objValList := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			item, _ := val.M__getitem__(gpy.Int(i))
			objValList[i] = PyObject2Value(item)
		}
		objVal = objValList
	case gpy.Tuple:
		objValList := make([]any, len(val))
		for i := 0; i < len(val); i++ {
			objValList[i] = PyObject2Value(val[i])
		}
		objVal = objValList
	case *gpy.Type:
		objVal = PyObject2Value(val.Dict)
	case gpy.NoneType:
		objVal = nil
	case nil:
		objVal = nil
	default:
		errors.NewAniError("不支持的类型").TryPanic()
	}
	return objVal
}

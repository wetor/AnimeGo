package python

import (
	"reflect"

	"github.com/go-python/gpython/py"

	"github.com/wetor/AnimeGo/pkg/errors"
)

func StructToObject(src any) py.Object {
	tempSrcValue := reflect.ValueOf(src)
	var srcValue reflect.Value

	// *struct to struct
	if tempSrcValue.Type().Kind() == reflect.Pointer {
		if tempSrcValue.IsNil() {
			return nil
		}
		srcValue = tempSrcValue.Elem()
	} else {
		srcValue = tempSrcValue
	}

	srcType := srcValue.Type()
	dst := py.NewStringDictSized(srcType.NumField())
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		value := srcValue.Field(i).Interface()

		keyName := field.Tag.Get("json")
		if len(keyName) == 0 {
			keyName = field.Name
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Pointer:
			dst[keyName] = StructToObject(value)
		default:
			dst[keyName] = ToObject(value)
		}
	}
	return dst
}

func ToObject(goVal any) py.Object {
	var pyObj py.Object
	switch val := goVal.(type) {
	case nil:
		pyObj = py.None
	case bool:
		pyObj = py.NewBool(val)
	case int8:
		pyObj = py.Int(val)
	case int:
		pyObj = py.Int(val)
	case int32:
		pyObj = py.Int(val)
	case int64:
		pyObj = py.Int(val)
	case float32:
		pyObj = py.Float(val)
	case float64:
		pyObj = py.Float(val)
	case string:
		pyObj = py.String(val)
	case map[string]any:
		pyValDict := py.NewStringDictSized(len(val))
		for key, value := range val {
			pyValDict[key] = ToObject(value)
		}
		pyObj = pyValDict
	default:
		refVal := reflect.ValueOf(goVal)
		switch refVal.Kind() {
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			l := refVal.Len()
			pyValList := py.NewListWithCapacity(l)
			for i := 0; i < l; i++ {
				pyValList.Append(ToObject(refVal.Index(i).Interface()))
			}
			pyObj = pyValList
		case reflect.Struct:
			fallthrough
		case reflect.Pointer:
			pyObj = StructToObject(goVal)
		default:
			errors.NewAniErrorf("不支持的类型: %v ", reflect.ValueOf(goVal).Type()).TryPanic()
		}
	}
	return pyObj
}

func ToValue(pyObj py.Object) any {
	var goVal any
	switch val := pyObj.(type) {
	case nil:
		goVal = nil
	case py.NoneType:
		goVal = nil
	case py.Bool:
		goVal = bool(val)
	case py.Int:
		goVal = int64(val)
	case py.Float:
		goVal = float64(val)
	case py.String:
		goVal = string(val)
	case py.StringDict:
		objValDict := make(map[string]any, len(val))
		for key, value := range val {
			objValDict[key] = ToValue(value)
		}
		goVal = objValDict
	case py.Tuple:
		objValList := make([]any, len(val))
		for i := 0; i < len(val); i++ {
			objValList[i] = ToValue(val[i])
		}
		goVal = objValList
	case *py.List:
		objValList := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			objValList[i] = ToValue(val.Items[i])
		}
		goVal = objValList
	case *py.Type:
		goVal = ToValue(val.Dict)
	default:
		errors.NewAniErrorf("不支持的类型: %v ", reflect.ValueOf(pyObj).Type()).TryPanic()
	}
	return goVal
}

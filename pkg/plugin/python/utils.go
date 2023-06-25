package python

import (
	"reflect"

	"github.com/pkg/errors"

	"github.com/go-python/gpython/py"
	"github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
)

func StructToObject(src any) (py.Object, error) {
	tempSrcValue := reflect.ValueOf(src)
	var srcValue reflect.Value

	// *struct to struct
	if tempSrcValue.Type().Kind() == reflect.Pointer {
		if tempSrcValue.IsNil() {
			return py.None, nil
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
			obj, err := StructToObject(value)
			if err != nil {
				return nil, err
			}
			dst[keyName] = obj
		default:
			obj, err := ToObject(value)
			if err != nil {
				return nil, err
			}
			dst[keyName] = obj
		}
	}
	return dst, nil
}

func ToObject(goVal any) (py.Object, error) {
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
			obj, err := ToObject(value)
			if err != nil {
				return nil, err
			}
			pyValDict[key] = obj
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
				obj, err := ToObject(refVal.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				pyValList.Append(obj)
			}
			pyObj = pyValList
		case reflect.Struct:
			fallthrough
		case reflect.Pointer:
			obj, err := StructToObject(goVal)
			if err != nil {
				return nil, err
			}
			pyObj = obj
		default:
			err := errors.WithStack(exceptions.ErrPluginTypeNotSupported{Type: reflect.ValueOf(goVal).Type()})
			log.DebugErr(err)
			return nil, err
		}
	}
	return pyObj, nil
}

func ToValue(pyObj py.Object) (any, error) {
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
			obj, err := ToValue(value)
			if err != nil {
				return nil, err
			}
			objValDict[key] = obj
		}
		goVal = objValDict
	case py.Tuple:
		objValList := make([]any, len(val))
		for i := 0; i < len(val); i++ {
			obj, err := ToValue(val[i])
			if err != nil {
				return nil, err
			}
			objValList[i] = obj
		}
		goVal = objValList
	case *py.List:
		objValList := make([]any, val.Len())
		for i := 0; i < val.Len(); i++ {
			obj, err := ToValue(val.Items[i])
			if err != nil {
				return nil, err
			}
			objValList[i] = obj
		}
		goVal = objValList
	case *py.Type:
		obj, err := ToValue(val.Dict)
		if err != nil {
			return nil, err
		}
		goVal = obj
	default:
		err := errors.WithStack(exceptions.ErrPluginTypeNotSupported{Type: reflect.ValueOf(pyObj).Type()})
		log.DebugErr(err)
		return nil, err
	}
	return goVal, nil
}

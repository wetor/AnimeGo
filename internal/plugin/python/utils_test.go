package python

import (
	"fmt"
	"github.com/go-python/gpython/py"
	"github.com/wetor/AnimeGo/internal/models"
	"testing"
)

func TestPyObject2Object(t *testing.T) {
	dict := py.NewStringDict()
	dict.M__setitem__(py.String("key1"), py.NewListFromStrings([]string{
		"list value1", "list_value2",
	}))
	dict2 := py.NewStringDict()
	dict2.M__setitem__(py.String("key_key1"), py.Int(132123132))
	dict.M__setitem__(py.String("dict1"), dict2)
	dict.M__setitem__(py.String("tuple"), py.Tuple{
		py.String("1111"),
		py.String("2222"),
	})
	fmt.Println(dict)

	obj := PyObject2Value(dict)
	fmt.Println(obj)
}

func TestValue2PyObject(t *testing.T) {
	object := models.Object{
		"key1": 10086,
		"key2": true,
		"keyObject": models.Object{
			"objKey1": "test123",
		},
		"keyObjList": []models.Object{
			{
				"objKey1": "test123",
			},
			{
				"objKey2": "6666",
			},
		},
		"keyStrList": []string{
			"这是文本1",
			"strList2",
		},
		"keyBoolList": []bool{
			true,
			false,
		},
	}
	pyObj := Value2PyObject(object)
	fmt.Println(pyObj)
}

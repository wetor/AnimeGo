package python_test

import (
	"fmt"
	"testing"

	"github.com/go-python/gpython/py"

	"github.com/wetor/AnimeGo/pkg/plugin/python"
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

	obj := python.ToValue(dict)
	fmt.Println(obj)
}

func TestValue2PyObject(t *testing.T) {
	object := map[string]any{
		"key1": 10086,
		"key2": true,
		"keyObject": map[string]any{
			"objKey1": "test123",
		},
		"keyObjList": []map[string]any{
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
	pyObj := python.ToObject(object)
	fmt.Println(pyObj)
}

func TestStructToObject(t *testing.T) {
	type InnerStruct struct {
		Best    int    `json:"in_best"`
		TestKey string `json:"in_test_key"`
	}
	type Struct struct {
		Best    int         `json:"best"`
		TestKey string      `json:"test_key"`
		Inner   InnerStruct `json:"inner"`
	}
	model := Struct{
		Best:    666,
		TestKey: "这是",
		Inner: InnerStruct{
			Best:    777,
			TestKey: "我",
		},
	}
	obj := python.ToObject(model)
	fmt.Println(obj)
}

package utils_test

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/utils"
	"reflect"
	"testing"
)

func TestSettings_Tag(t *testing.T) {
	tag := "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}"
	got := utils.Tag(tag, "2022-04-11", 10)
	want := "2022年4月新番,2,春季,第10集,周1,星期一"
	if got != want {
		t.Errorf("Tag() = %v, want %v", got, want)
	}
}

func TestMapToStruct(t *testing.T) {
	obj := models.Object{
		"best":    100,
		"testKey": "这是字符串",
	}
	type Struct struct {
		Best    int64  `json:"best"`
		TestKey string `json:"testKey"`
	}
	model := Struct{}
	utils.MapToStruct(obj, &model)
	fmt.Println(model)
}

func TestStructToMap(t *testing.T) {
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
	obj := utils.StructToMap(&model)
	fmt.Println(obj)
	fmt.Println(obj["best"], reflect.ValueOf(obj["best"]).Type())
}

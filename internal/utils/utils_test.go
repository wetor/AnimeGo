package utils_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
)

func TestSettings_Tag(t *testing.T) {
	tag := "{year}年{quarter}月新番,{quarter_index},{quarter_name}季,第{ep}集,周{week},{week_name}"
	got := utils.Tag(tag, "2022-04-11", 10)
	want := "2022年4月新番,2,春季,第10集,周1,星期一"
	if got != want {
		t.Errorf("Tag() = %v, want %v", got, want)
	}
}

func TestMap2ModelByJson(t *testing.T) {
	obj := models.Object{
		"best":    100,
		"testKey": "这是字符串",
	}
	type Struct struct {
		Best    int64  `json:"best"`
		TestKey string `json:"testKey"`
	}
	model := Struct{}
	utils.Map2ModelByJson(obj, &model)
	fmt.Println(model)
}

func TestMap2Model(t *testing.T) {
	obj := models.Object{
		"Best":    100,
		"TestKey": "这是字符串",
	}
	type Struct struct {
		Best    int
		TestKey string
	}
	model := Struct{}
	utils.Map2Model(obj, &model)
	fmt.Println(model)
}

func TestModel2Map(t *testing.T) {
	obj := models.Object{}
	type Struct struct {
		Best    int
		TestKey string
	}
	model := Struct{
		Best:    666,
		TestKey: "这是",
	}
	utils.Model2Map(&model, obj)
	fmt.Println(obj)
}

func TestModel2MapByJson(t *testing.T) {
	obj := models.Object{}
	type Struct struct {
		Best    int    `json:"best"`
		TestKey string `json:"testKey"`
	}
	model := Struct{
		Best:    666,
		TestKey: "这是",
	}
	utils.Model2MapByJson(&model, obj)
	fmt.Println(obj)
	fmt.Println(obj["best"], reflect.ValueOf(obj["best"]).Type())
}

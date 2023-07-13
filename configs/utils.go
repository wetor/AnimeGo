package configs

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/models"
)

func ConvertPluginInfo(info []PluginInfo) []models.Plugin {
	plugins := make([]models.Plugin, len(info))
	_ = copier.Copy(&plugins, &info)
	return plugins
}

func getFieldByTag(value reflect.Value, tag string) reflect.Value {
	tags := strings.Split(tag, ".")
	if len(tags) == 1 {
		return value.FieldByName(tag)
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := value.Type().Field(i)

		if fieldType.Name == tags[0] {
			if field.Kind() == reflect.Struct {
				return getFieldByTag(field, strings.Join(tags[1:], "."))
			}
			break
		}
	}

	return reflect.Value{}
}

func Env2Config(env *Environment, conf *Config, prefix string) error {
	envValue := reflect.ValueOf(env).Elem()
	confValue := reflect.ValueOf(conf).Elem()

	for i := 0; i < envValue.NumField(); i++ {
		envField := envValue.Field(i)
		confField := getFieldByTag(confValue, envValue.Type().Field(i).Tag.Get("val"))

		if !confField.IsValid() {
			return errors.New(fmt.Sprintf("Field %s not found in Config", envValue.Type().Field(i).Tag.Get("val")))
		}

		if envField.IsNil() {
			continue
		}
		envTag := envValue.Type().Field(i).Tag.Get("env")
		envVar := envField.Elem()
		log.Printf("发现环境变量 %s%s=%v", prefix, envTag, envVar)
		switch envField.Elem().Kind() {
		case reflect.String:
			confField.SetString(envVar.Interface().(string))
		case reflect.Int:
			confField.SetInt(int64(envVar.Interface().(int)))
		default:
			return errors.New(fmt.Sprintf("Field %s has an unsupported type", envTag))
		}
	}

	return nil
}

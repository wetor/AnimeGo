package configs

import (
	"fmt"
	"log"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	encoder "github.com/wetor/AnimeGo/third_party/yaml-encoder"
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

func Env2Config(env *Environment, conf *Config, prefix string) (bool, error) {
	rewrite := false
	envValue := reflect.ValueOf(env).Elem()
	confValue := reflect.ValueOf(conf).Elem()

	for i := 0; i < envValue.NumField(); i++ {
		envField := envValue.Field(i)
		confField := getFieldByTag(confValue, envValue.Type().Field(i).Tag.Get("val"))

		if !confField.IsValid() {
			return rewrite, errors.New(fmt.Sprintf("Field %s not found in Config", envValue.Type().Field(i).Tag.Get("val")))
		}

		if envField.IsNil() {
			continue
		}
		envTag := envValue.Type().Field(i).Tag.Get("env")
		envVar := envField.Elem()
		log.Printf("发现环境变量 %s%s=%v", prefix, envTag, envVar)
		rewrite = true
		switch envField.Elem().Kind() {
		case reflect.String:
			confField.SetString(envVar.Interface().(string))
		case reflect.Int:
			confField.SetInt(int64(envVar.Interface().(int)))
		default:
			return rewrite, errors.New(fmt.Sprintf("Field %s has an unsupported type", envTag))
		}
	}

	return rewrite, nil
}

func encodeConfig(conf any) ([]byte, error) {
	defaultSettingComment()
	defaultAdvancedComment()
	yml := encoder.NewEncoder(conf,
		encoder.WithComments(encoder.CommentsOnHead),
		encoder.WithCommentsMap(configComment),
	)
	content, err := yml.Encode()
	if err != nil {
		return nil, err
	}
	return content, nil
}

func BackupConfig(file string, version string) error {
	dir, name := path.Split(file)
	ext := path.Ext(name)
	name = strings.TrimSuffix(name, ext)
	timeStr := time.Now().Format("20060102150405")
	name = fmt.Sprintf("%s-%s-%s%s", name, version, timeStr, ext)
	oldFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	out := path.Join(dir, name)
	err = os.WriteFile(out, oldFile, constant.WriteFilePerm)
	if err != nil {
		return err
	}
	log.Printf("备份原配置文件到：'%s'\n", out)
	return nil
}

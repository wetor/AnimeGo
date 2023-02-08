package utils

import (
	"strings"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
)

type Schema struct {
	Name     string
	Optional bool
}

func ParseSchemas(schemas []string) []*Schema {
	paramsSchema := make([]*Schema, len(schemas))
	for i, schema := range schemas {
		s := strings.Split(schema, ",")
		paramsSchema[i] = &Schema{
			Name:     s[0],
			Optional: false,
		}
		if len(s) > 1 && s[1] == "optional" {
			paramsSchema[i].Optional = true
		}
	}
	return paramsSchema
}

func CheckSchema(schemas []*Schema, object any) {
	objectMap, ok := object.(models.Object)
	if !ok {
		errors.NewAniError("类型错误").TryPanic()
	}

	for _, schema := range schemas {
		if !schema.Optional {
			has := false
			for key := range objectMap {
				if key == schema.Name {
					has = true
					break
				}
			}
			if !has {
				errors.NewAniError("缺少参数: " + schema.Name).TryPanic()
			}
		}
	}

	for key := range objectMap {
		has := false
		for _, schema := range schemas {
			if key == schema.Name {
				has = true
				break
			}
		}
		if !has {
			errors.NewAniError("多余参数: " + key).TryPanic()
		}
	}
}

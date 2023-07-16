package plugin

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/pkg/exceptions"
	"github.com/wetor/AnimeGo/pkg/log"
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

func CheckSchema(schemas []*Schema, object any) error {
	objectMap := object.(map[string]any)

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
				err := errors.WithStack(&exceptions.ErrPluginSchemaMissing{Name: schema.Name})
				log.DebugErr(err)
				return err
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
			err := errors.WithStack(&exceptions.ErrPluginSchemaUnknown{Name: key})
			log.DebugErr(err)
			return err
		}
	}
	return nil
}

package components

import (
	"fmt"
	g "github.com/maragudk/gomponents"
	"reflect"
	"strings"
)

type NameInfo struct {
	Name        string
	DisplayName string
	Comment     string
	Options     []string
	Hidden      bool
}

func (n NameInfo) ID(t string) string {
	return t + "-" + n.Name
}

func (n NameInfo) Last() string {
	return n.Name[strings.LastIndex(n.Name, "-")+1:]
}

type Components struct {
	commentMap map[string]string
}

func NewComponents(commentMap map[string]string) *Components {
	return &Components{
		commentMap: commentMap,
	}
}

func (c Components) comment2Html(comment string) string {
	html := strings.Builder{}
	lines := strings.Split(comment, "\n")
	isLi := false
	for i, line := range lines {
		if !isLi && len(line) > 2 && line[:2] == "  " {
			isLi = true
			html.WriteString("<ul>")
		}
		if isLi {
			html.WriteString("<li>")
			html.WriteString(line)
			html.WriteString("</li>")
		} else if i < len(lines)-1 {
			html.WriteString(line)
			html.WriteString("<br>")
		} else {
			html.WriteString(line)
		}
	}
	if isLi {
		html.WriteString("</ul>")
	}
	return html.String()
}

func (c Components) getFiledName(name string, t reflect.StructField) NameInfo {
	fieldName := t.Tag.Get("json")
	if len(fieldName) == 0 {
		fieldName = t.Name
	}
	displayName := t.Tag.Get("attr")
	if len(displayName) == 0 {
		displayName = fieldName
	}

	comment := t.Tag.Get("comment")
	if len(comment) == 0 {
		comment = displayName
		commentKey := t.Tag.Get("comment_key")
		if len(commentKey) == 0 || c.commentMap == nil {
			comment = displayName
		} else if tmp, ok := c.commentMap[commentKey]; ok {
			comment = tmp
		}
	}

	if len(comment) == 0 {
		comment = displayName
	}

	opts := strings.Split(t.Tag.Get("options"), ",")

	comment = c.comment2Html(comment)
	return NameInfo{
		Name:        name + "-" + fieldName,
		DisplayName: displayName,
		Comment:     comment,
		Options:     opts,
		Hidden:      false,
	}
}

func (c Components) ArrayAdder(name NameInfo, t reflect.Type) g.Node {
	typeObject := t.Elem()
	ptr := reflect.New(typeObject)
	valueObject := ptr.Elem()
	name.Hidden = true
	return ArrayItem(name, -1, c.Struct2Node(name, valueObject.Interface())...)

}

func (c Components) Struct2Node(name NameInfo, object any) []g.Node {
	tempObject := reflect.ValueOf(object)
	var valueObject reflect.Value

	if tempObject.Type().Kind() == reflect.Pointer {
		if tempObject.IsNil() {
			return nil
		}
		valueObject = tempObject.Elem()
	} else {
		valueObject = tempObject
	}

	result := make([]g.Node, 0)
	typeObject := valueObject.Type()
	if typeObject.Kind() != reflect.Struct {
		return []g.Node{c.Value2Node(name, object)}
	}

	for i := 0; i < typeObject.NumField(); i++ {
		typeField := typeObject.Field(i)
		fieldName := c.getFiledName(name.Name, typeField)
		fieldName.Hidden = name.Hidden

		valueField := valueObject.Field(i)
		value := valueField.Interface()
		switch typeField.Type.Kind() {
		case reflect.Pointer:
			fallthrough
		case reflect.Struct:
			result = append(result, StructCard(fieldName, c.Struct2Node(fieldName, value)...))
		case reflect.Array:
			fallthrough
		case reflect.Slice:
			array := make([]g.Node, valueField.Len())
			for j := 0; j < valueField.Len(); j++ {
				v := valueField.Index(j)
				elemName := NameInfo{
					Name:        fmt.Sprintf("%s-%d", fieldName.Name, j),
					DisplayName: fmt.Sprintf("[%d]", j),
					Hidden:      fieldName.Hidden,
				}
				array[j] = ArrayItem(elemName, j, c.Struct2Node(elemName, v.Interface())...)
			}
			array = append(array, c.ArrayAdder(fieldName, typeField.Type))
			result = append(result, ArrayList(fieldName, array...))
		default:
			result = append(result, c.Value2Node(fieldName, value))
		}

	}
	return result
}

func (c Components) Value2Node(name NameInfo, object any) g.Node {
	switch val := object.(type) {
	case string:
		if len(name.Options) == 1 && name.Options[0] == "" {
			return StringInput(InputOptions{
				NameInfo: name,
				Value:    val,
			})
		} else {
			m := make(map[string]string, len(name.Options))
			for _, opt := range name.Options {
				m[opt] = opt
			}
			return SelectInput(name, val, m)
		}
	case int:
		return NumberInput(InputOptions{
			NameInfo: name,
			Value:    val,
		})
	case bool:
		return BoolInput(InputOptions{
			NameInfo: name,
			Value:    val,
		})
	}
	return nil
}

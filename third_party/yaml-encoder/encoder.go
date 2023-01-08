package yaml_encoder

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Encoder implements config encoder.
type Encoder struct {
	value   interface{}
	options *Options
}

// CommentsFlag comments encoding flags type.
type CommentsFlag int

func (f CommentsFlag) enabled() bool {
	return f > CommentsDisabled
}

const (
	CommentsDisabled CommentsFlag = iota
	CommentsOnHead
	CommentsInLine
	CommentsOnFoot
)

const (
	AttrCommentSep = ". "
)

// Options defines encoder config.
type Options struct {
	CommentsFlag CommentsFlag
	OmitEmpty    bool
	CommentsMap  map[string]string
}

func newOptions(opts ...Option) *Options {
	res := &Options{
		CommentsFlag: CommentsDisabled,
		OmitEmpty:    true,
	}

	for _, o := range opts {
		o(res)
	}

	return res
}

// Option gives ability to alter config encoder output settings.
type Option func(*Options)

// WithComments enables comments in the encoder.
func WithComments(flag CommentsFlag) Option {
	return func(o *Options) {
		o.CommentsFlag = flag
	}
}

// WithOmitEmpty toggles omitempty handling.
func WithOmitEmpty(value bool) Option {
	return func(o *Options) {
		o.OmitEmpty = value
	}
}

// WithCommentsMap enables comments from map.
func WithCommentsMap(m map[string]string) Option {
	return func(o *Options) {
		o.CommentsMap = m
	}
}

// NewEncoder initializes and returns an `Encoder`.
func NewEncoder(value interface{}, opts ...Option) *Encoder {
	return &Encoder{
		value:   value,
		options: newOptions(opts...),
	}
}

// Marshal converts value to YAML-serializable value (suitable for MarshalYAML).
func (e *Encoder) Marshal() (*yaml.Node, error) {
	node, err := toYamlNode(e.value, e.options)
	if err != nil {
		return nil, err
	}
	return node, nil
}

// Encode converts value to yaml.
//nolint:gocyclo
func (e *Encoder) Encode() ([]byte, error) {
	if e.options.CommentsFlag == CommentsDisabled {
		return yaml.Marshal(e.value)
	}

	node, err := e.Marshal()
	if err != nil {
		return nil, err
	}

	// special handling for case when we get an empty output
	if node.Kind == yaml.MappingNode && len(node.Content) == 0 && node.FootComment != "" && e.options.CommentsFlag.enabled() {
		res := ""

		if node.HeadComment != "" {
			res += node.HeadComment + "\n"
		}

		lines := strings.Split(res+node.FootComment, "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
				continue
			}

			lines[i] = "# " + line
		}

		return []byte(strings.Join(lines, "\n")), nil
	}
	return yaml.Marshal(node)
}

// EncodeDoc converts comment to json.
func (e *Encoder) EncodeDoc() ([]byte, error) {
	if e.options.CommentsFlag == CommentsDisabled {
		return yaml.Marshal(e.value)
	}
	node, err := e.Marshal()
	if err != nil {
		return nil, err
	}
	jsonStr := e.encodeDoc(node)

	return []byte(jsonStr[:len(jsonStr)-1]), nil
}

var (
	replacer = strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`)
	spliter  = func(s string) (attr, comment string) {
		i := strings.Index(s, AttrCommentSep)
		if i < 0 {
			return s, ""
		}
		return s[:i], s[i+len(AttrCommentSep):]
	}
)

func (e *Encoder) encodeDoc(node *yaml.Node) string {
	buf := bytes.NewBuffer(nil)
	if node.Kind == yaml.MappingNode && len(node.Content) > 0 {
		buf.WriteString("{")
		if len(node.Value) > 0 {
			attr, comment := spliter(node.Value)
			buf.WriteString(fmt.Sprintf(`"_attr":"%s",`, attr))
			if len(comment) > 0 {
				buf.WriteString(fmt.Sprintf(`"_comment":"%s",`, replacer.Replace(comment)))
			}
		}
		isKey := true
		comment := ""
		for _, n := range node.Content {
			if isKey {
				comment = n.HeadComment
				buf.WriteString(fmt.Sprintf(`"%s":`, n.Value))
				isKey = false
			} else {
				n.Value = comment
				buf.WriteString(e.encodeDoc(n))
				isKey = true
			}
		}
		if buf.Len() > 1 {
			buf.Truncate(buf.Len() - 1)
		}
		buf.WriteString("},")

	} else if node.Kind == yaml.ScalarNode || node.Kind == yaml.SequenceNode {
		attr, comment := spliter(node.Value)
		buf.WriteString("{")
		buf.WriteString(fmt.Sprintf(`"_attr":"%s"`, attr))
		if len(comment) > 0 {
			buf.WriteString(fmt.Sprintf(`,"_comment":"%s"`, replacer.Replace(comment)))
		}
		buf.WriteString("},")
	}
	return buf.String()
}

func isEmpty(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Ptr:
		return value.IsNil()
	case reflect.Map:
		return len(value.MapKeys()) == 0
	case reflect.Slice:
		return value.Len() == 0
	default:
		return value.IsZero()
	}
}

func isNil(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	//nolint:exhaustive
	switch value.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

//nolint:gocyclo,cyclop
func toYamlNode(in interface{}, options *Options) (*yaml.Node, error) {
	node := &yaml.Node{}

	// flags := options.Comments

	// do not wrap yaml.Node into yaml.Node
	if n, ok := in.(*yaml.Node); ok {
		return n, nil
	}

	// if input implements yaml.Marshaler we should use that marshaller instead
	// same way as regular yaml marshal does
	if m, ok := in.(yaml.Marshaler); ok && !isNil(reflect.ValueOf(in)) {
		res, err := m.MarshalYAML()
		if err != nil {
			return nil, err
		}

		if n, ok := res.(*yaml.Node); ok {
			return n, nil
		}

		in = res
	}

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	//nolint:exhaustive
	switch v.Kind() {
	case reflect.Struct:
		node.Kind = yaml.MappingNode

		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			// skip unexported fields
			if !v.Field(i).CanInterface() {
				continue
			}

			tag := t.Field(i).Tag.Get("yaml")
			comment, has := t.Field(i).Tag.Lookup("comment")
			if options.CommentsMap != nil && !has {
				commentKey, has := t.Field(i).Tag.Lookup("comment_key")
				// default use yaml
				if !has {
					commentKey = tag
				}
				comment = options.CommentsMap[commentKey]
			}
			attr := t.Field(i).Tag.Get("attr")

			parts := strings.Split(tag, ",")
			fieldName := parts[0]
			parts = parts[1:]

			if fieldName == "" {
				fieldName = strings.ToLower(t.Field(i).Name)
			}

			if fieldName == "-" {
				continue
			}

			var (
				empty  = isEmpty(v.Field(i))
				skip   bool
				inline bool
				flow   bool
			)

			for _, part := range parts {
				if part == "omitempty" && empty && options.OmitEmpty {
					skip = true
				}

				if part == "inline" {
					inline = true
				}

				if part == "flow" {
					flow = true
				}
			}

			var value interface{}
			if v.Field(i).CanInterface() {
				value = v.Field(i).Interface()
			}

			if skip {
				continue
			}

			var style yaml.Style
			if flow {
				style |= yaml.FlowStyle
			}
			if inline {
				child, err := toYamlNode(value, options)
				if err != nil {
					return nil, err
				}

				if child.Kind == yaml.MappingNode || child.Kind == yaml.SequenceNode {
					appendNodes(node, child.Content...)
				}
			} else if err := addToMap(node, fieldName, value, comment, attr, style, options); err != nil {
				return nil, err
			}
		}
	case reflect.Map:
		node.Kind = yaml.MappingNode
		keys := v.MapKeys()
		// always interate keys in alphabetical order to preserve the same output for maps
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			element := v.MapIndex(k)
			value := element.Interface()

			if err := addToMap(node, k.Interface(), value, "", "", 0, options); err != nil {
				return nil, err
			}
		}
	case reflect.Slice:
		node.Kind = yaml.SequenceNode
		nodes := make([]*yaml.Node, v.Len())

		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)

			var err error

			nodes[i], err = toYamlNode(element.Interface(), options)
			if err != nil {
				return nil, err
			}
		}
		appendNodes(node, nodes...)
	default:
		if err := node.Encode(in); err != nil {
			return nil, err
		}
	}

	return node, nil
}

func appendNodes(dest *yaml.Node, nodes ...*yaml.Node) {
	if dest.Content == nil {
		dest.Content = []*yaml.Node{}
	}

	dest.Content = append(dest.Content, nodes...)
}

func addToMap(dest *yaml.Node, fieldName, in interface{}, fieldComment, fieldAttr string, style yaml.Style, options *Options) error {
	key, err := toYamlNode(fieldName, options)
	if err != nil {
		return err
	}

	value, err := toYamlNode(in, options)
	if err != nil {
		return err
	}
	value.Style = style

	if options.CommentsFlag.enabled() {
		addComment(key, fieldComment, fieldAttr, options.CommentsFlag)
	}
	appendNodes(dest, key, value)

	return nil
}

func addComment(node *yaml.Node, comment, attr string, flag CommentsFlag) {
	if flag.enabled() {
		dest := []*string{
			&node.HeadComment,
			&node.LineComment,
			&node.FootComment,
		}
		*dest[int(flag)-1] = comment
		if len(attr) > 0 {
			if len(*dest[int(CommentsOnHead)-1]) == 0 {
				*dest[int(CommentsOnHead)-1] = attr
			} else {
				*dest[int(CommentsOnHead)-1] = attr + AttrCommentSep + *dest[int(CommentsOnHead)-1]
			}
		}
	}
}

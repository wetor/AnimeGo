package plugin

type Plugin interface {
	Execute(file string, params Object) any
	SetSchema(paramsSchema, resultSchema []string)
}

// Object 对象类型
type Object map[string]any

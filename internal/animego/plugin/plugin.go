package plugin

type Plugin interface {
	Execute(file string, params map[string]any) any
}

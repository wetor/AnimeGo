package plugin

type Plugin interface {
	Execute(file string, params map[string]interface{}) (interface{}, error)
}

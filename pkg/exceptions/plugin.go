package exceptions

import "fmt"

type ErrPlugin struct {
	Type    string
	File    string
	Message any
}

func (e ErrPlugin) Error() string {
	return fmt.Sprintf("[%s] %s: %v", e.Type, e.File, e.Message)
}

type ErrPluginSchemaMissing struct {
	Name string
}

func (e ErrPluginSchemaMissing) Error() string {
	return fmt.Sprintf("缺少参数: %s", e.Name)
}

type ErrPluginSchemaUnknown struct {
	Name string
}

func (e ErrPluginSchemaUnknown) Error() string {
	return fmt.Sprintf("未知参数: %s", e.Name)
}

type ErrPluginTypeNotSupported struct {
	Type any
}

func (e ErrPluginTypeNotSupported) Error() string {
	return fmt.Sprintf("不支持的类型: %v", e.Type)
}

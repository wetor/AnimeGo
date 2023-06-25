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

type ErrPluginDisabled struct {
	Type string
	File string
}

func (e ErrPluginDisabled) Error() string {
	return fmt.Sprintf("[%s] %s: 插件未启用", e.Type, e.File)
}

type ErrParseFailed struct {
}

func (e ErrParseFailed) Error() string {
	return "解析季度失败，未设置默认值，结束此流程"
}

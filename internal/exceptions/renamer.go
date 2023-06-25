package exceptions

import "fmt"

type ErrRename struct {
	Src     string
	Message string
}

func (e ErrRename) Error() string {
	return fmt.Sprintf("重命名 %s 失败: %s", e.Src, e.Message)
}

type ErrRenameStep struct {
	Src     string
	Step    string
	Message string
}

func (e ErrRenameStep) Error() string {
	return fmt.Sprintf("%s %s 失败: %s", e.Step, e.Src, e.Message)
}

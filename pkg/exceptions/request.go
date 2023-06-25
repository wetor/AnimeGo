package exceptions

import "fmt"

type ErrRequest struct {
	Name string
}

func (e ErrRequest) Error() string {
	return fmt.Sprintf("请求 %s 失败", e.Name)
}

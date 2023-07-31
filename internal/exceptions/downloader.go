package exceptions

import "fmt"

type ErrClientExistItem struct {
	Client string
	Name   string
}

func (e ErrClientExistItem) Error() string {
	return fmt.Sprintf("%s 正在下载: %s", e.Client, e.Name)
}

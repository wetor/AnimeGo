package exceptions

import "fmt"

type ErrClientExistItem struct {
	Client string
	Name   string
}

func (e ErrClientExistItem) Error() string {
	return fmt.Sprintf("%s 正在下载: %s", e.Client, e.Name)
}

func (e ErrClientExistItem) Exist() bool {
	return true
}

type ErrDownloadExist struct {
	Name string
}

func (e ErrDownloadExist) Error() string {
	return fmt.Sprintf("已下载: %s", e.Name)
}

func (e ErrDownloadExist) Exist() bool {
	return true
}

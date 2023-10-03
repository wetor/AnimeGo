package exceptions

import "fmt"

type ErrClient struct {
	Client  string
	Message string
}

func (e ErrClient) Error() string {
	return fmt.Sprintf("%s 下载器错误: %s", e.Client, e.Message)
}

type ErrClientNoConnected struct {
	Client string
}

func (e ErrClientNoConnected) Error() string {
	return fmt.Sprintf("未连接到 %s 下载器", e.Client)
}

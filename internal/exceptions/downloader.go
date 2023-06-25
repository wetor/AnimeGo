package exceptions

import "fmt"

type ErrDownloader struct {
	Client  string
	Message string
}

func (e ErrDownloader) Error() string {
	return fmt.Sprintf("%s 下载器错误: %s", e.Client, e.Message)
}

type ErrDownloaderNoConnected struct {
	Client string
}

func (e ErrDownloaderNoConnected) Error() string {
	return fmt.Sprintf("未连接到 %s 下载器", e.Client)
}

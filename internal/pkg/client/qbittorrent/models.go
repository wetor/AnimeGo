package qbittorrent

import "sync"

type Conf struct {
	Url                  string
	Username             string
	Password             string
	DownloadPath         string
	ConnectTimeoutSecond int // 连接超时
	CheckTimeSecond      int // 定时检查是否在线
	RetryConnectNum      int // 连接失败重试次数
}

type Options struct {
	Url          string
	Username     string
	Password     string
	DownloadPath string

	ConnectTimeoutSecond int // 连接超时
	CheckTimeSecond      int // 定时检查是否在线
	RetryConnectNum      int // 连接失败重试次数
	WG                   *sync.WaitGroup
}

func (o *Options) Default() {
	if o.ConnectTimeoutSecond == 0 {
		o.ConnectTimeoutSecond = 5
	}
	if o.RetryConnectNum == 0 {
		o.RetryConnectNum = 3
	}
	if o.CheckTimeSecond == 0 {
		o.CheckTimeSecond = 60
	}
}

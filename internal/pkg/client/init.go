package client

import (
	"context"
	"sync"
)

var (
	DownloadPath         string
	SeedingTimeMinute    int // 分钟
	ConnectTimeoutSecond int // 连接超时
	CheckTimeSecond      int // 定时检查是否在线
	RetryConnectNum      int // 连接失败重试次数

	WG  *sync.WaitGroup
	Ctx context.Context
)

type Options struct {
	DownloadPath         string
	SeedingTimeMinute    int
	ConnectTimeoutSecond int // 连接超时
	CheckTimeSecond      int // 定时检查是否在线
	RetryConnectNum      int // 连接失败重试次数
	WG                   *sync.WaitGroup
	Ctx                  context.Context
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

func Init(opt *Options) {
	opt.Default()
	SeedingTimeMinute = opt.SeedingTimeMinute
	ConnectTimeoutSecond = opt.ConnectTimeoutSecond
	RetryConnectNum = opt.RetryConnectNum
	CheckTimeSecond = opt.CheckTimeSecond
	DownloadPath = opt.DownloadPath

	WG = opt.WG
	Ctx = opt.Ctx
}

package models

import (
	"context"
	"sync"
)

type ClientOptions struct {
	Ctx context.Context
	WG  *sync.WaitGroup

	Url      string
	Username string
	Password string

	DownloadPath         string
	SeedingTimeMinute    int
	ConnectTimeoutSecond int // 连接超时
	CheckTimeSecond      int // 定时检查是否在线
	RetryConnectNum      int // 连接失败重试次数
}

type ListOptions struct {
	Status   string
	Category string
	Tag      string
}

type AddOptions struct {
	Url      string
	File     string // optional torrent file
	SavePath string
	Category string
	Tag      string
	Name     string // 保存名字
}

type DeleteOptions struct {
	Hash       []string
	DeleteFile bool
}

type TorrentItem struct {
	Hash  string `json:"hash"`
	State string `json:"state"`

	ContentPath string  `json:"content_path"` // 非必要
	Name        string  `json:"name"`         // 非必要
	Progress    float64 `json:"progress"`     // 非必要
}

type Config struct {
	ApiUrl       string `json:"api_url"`
	DownloadPath string `json:"download_path"`
}

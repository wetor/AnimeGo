package downloader

import (
	"github.com/wetor/AnimeGo/internal/constant"
	"sync"
)

type ItemState struct {
	Torrent constant.TorrentState
	Notify  NotifyState
	Name    string
	Info    any // 下载项信息
}

type Options struct {
	RefreshSecond          int
	Category               string
	Tag                    string
	AllowDuplicateDownload bool
	WG                     *sync.WaitGroup
}

package downloader

import (
	"github.com/wetor/AnimeGo/internal/constant"
)

type ItemState struct {
	Torrent constant.TorrentState
	Notify  NotifyState
	Name    string
	Info    any // 下载项信息
}

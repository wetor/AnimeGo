package downloader

import "github.com/wetor/AnimeGo/internal/pkg/client"

type ItemState struct {
	Torrent client.TorrentState
	Notify  NotifyState
	Name    string
	Info    any // 下载项信息
}

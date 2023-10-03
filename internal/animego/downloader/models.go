package downloader

import "github.com/wetor/AnimeGo/internal/models"

type ItemState struct {
	Torrent models.TorrentState
	Notify  NotifyState
	Name    string
	Info    any // 下载项信息
}

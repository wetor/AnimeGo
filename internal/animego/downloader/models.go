package downloader

import "github.com/wetor/AnimeGo/internal/models"

type ItemState struct {
	Torrent models.TorrentState
	Notify  NotifyState
	Info    any // 下载项信息
}

package models

import "github.com/wetor/AnimeGo/internal/constant"

type ItemState struct {
	Torrent constant.TorrentState
	Notify  constant.NotifyState
	Name    string
	Info    any // 下载项信息
}

type ClientEvent struct {
	Gid  string `json:"gid"`  // aira2
	Hash string `json:"hash"` // qBittorrent
}

package models

import (
	"path/filepath"
)

const (
	MagnetType  = "magnet"
	TorrentType = "torrent"
)

type FeedItem struct {
	Url        string       `json:"url"`      // Link，详情页连接，用于下一步解析番剧信息
	Name       string       `json:"name"`     // 种子名
	Date       string       `json:"date"`     // 发布日期
	Type       string       `json:"type"`     // 下载类型，[application/x-bittorrent]
	Download   string       `json:"download"` // 下载链接
	Length     int64        `json:"length"`   // 种子大小
	NameParsed *TitleParsed // 标题解析信息
}

// DownloadType
//  @Description: 下载链接类型，[torrent, magnet, unknown]
//  @receiver FeedItem
//  @return string
//
func (i FeedItem) DownloadType() string {
	if len(i.Download) > 8 {
		if i.Download[:8] == "magnet:?" {
			return MagnetType
		}

		if i.Download[len(i.Download)-8:] == ".torrent" {
			return TorrentType
		}
	}
	return ""
}

// Hash
//  @Description: torrent 类型的hash
//  @receiver FeedItem
//  @return string
//
func (i FeedItem) Hash() string {
	if i.DownloadType() == TorrentType {
		_, hash := filepath.Split(i.Download)
		if len(hash) >= 40 {
			return hash[:40]
		}
	}
	return "ItemUrl_" + Md5Str(i.Url)
}

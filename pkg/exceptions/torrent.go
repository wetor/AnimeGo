package exceptions

import "fmt"

type ErrTorrent struct {
	Message string
}

func (e ErrTorrent) Error() string {
	return e.Message
}

type ErrTorrentLoad struct {
	Message string
}

func (e ErrTorrentLoad) Error() string {
	return fmt.Sprintf("加载Torrent失败: %s", e.Message)
}

type ErrTorrentParse struct {
	Message string
}

func (e ErrTorrentParse) Error() string {
	return fmt.Sprintf("解析Torrent失败: %s", e.Message)
}

type ErrTorrentHash struct {
	Hash    string
	Message string
}

func (e ErrTorrentHash) Error() string {
	return fmt.Sprintf("%s: %s", e.Message, e.Hash)
}

type ErrTorrentUrl struct {
	Url string
}

func (e ErrTorrentUrl) Error() string {
	return fmt.Sprintf("无法识别Torrent Url: %s", e.Url)
}

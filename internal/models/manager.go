package models

type TorrentState string

type Torrent struct {
	Hash       string
	State      TorrentState // 状态
	OldPath    string       // 下载到的旧路径
	Path       string       // 当前内容路径
	Renamed    bool         // 是否已重命名/移动
	Downloaded bool         // 是否已下载完成
	Scraped    bool         // 是否已经完成搜刮
}

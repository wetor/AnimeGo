package models

type TorrentState string

type Torrent struct {
	Hash       string
	State      TorrentState // 状态
	Path       string       // 当前内容路径
	Renamed    bool         // 是否已重命名，由下载器操作
	Downloaded bool         // 是否已下载完成
	Moved      bool         // 是否已经移动到正确位置，由程序操作
	Scraped    bool         // 是否已经完成搜刮
}

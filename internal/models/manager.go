package models

type TorrentState string

type Torrent struct {
	Hash       string
	State      TorrentState      // 状态
	OldPath    string            // 下载到的旧路径
	Path       string            // 当前内容路径
	Init       bool              // 是否初始化
	Renamed    bool              // 是否已重命名/移动
	StateChan  chan TorrentState // 通知chan
	Downloaded bool              // 是否已下载完成
	Seeded     bool              // 是否做种
	Scraped    bool              // 是否已经完成搜刮
}

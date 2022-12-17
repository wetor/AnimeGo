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

type DownloadStatus struct {
	Hash     string       `json:"hash"`
	State    TorrentState `json:"state"`
	Path     string       `json:"path"`      // 文件存储相对
	ExpireAt int64        `json:"expire_at"` // 过期时间

	Init       bool `json:"init"`       // 是否初始化
	Renamed    bool `json:"renamed"`    // 是否已重命名/移动
	Downloaded bool `json:"downloaded"` // 是否已下载完成
	Seeded     bool `json:"seeded"`     // 是否做种
	Scraped    bool `json:"scraped"`    // 是否已经完成搜刮
}

package models

type TorrentState struct {
	Hash    string
	State   string // 状态
	Path    string // 当前内容路径
	Renamed bool   // 是否已经重命名
}

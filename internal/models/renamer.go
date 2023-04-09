package models

type RenamePath struct {
	File []string
	Path string
}

type RenameResult struct {
	Index     int    `json:"index"`
	Filepath  string `json:"filepath"`
	TVShowDir string `json:"tvshow_dir"`
}

type RenameCallback func(*RenameResult)
type CompleteCallback func(*RenameResult)

type RenameOptions struct {
	Entity           *AnimeEntity
	SrcDir           string // 原名
	DstDir           string // 目标名
	Mode             string
	State            []chan TorrentState
	RenameCallback   RenameCallback   // 重命名完成后回调
	CompleteCallback CompleteCallback // 完成重命名所有流程后回调
}

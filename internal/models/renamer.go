package models

type RenameDst struct {
	Anime    *AnimeEntity
	Content  *TorrentContentItem
	SavePath string
}

type RenameResult struct {
	Filepath  string `json:"filepath"`
	TVShowDir string `json:"tvshow_dir"`
}

type RenameCallback func(*RenameResult)
type CompleteCallback func(*RenameResult)

type RenameOptions struct {
	Src              string     // 原名
	Dst              *RenameDst // 目标名
	Mode             string
	State            <-chan TorrentState
	RenameCallback   RenameCallback   // 重命名完成后回调
	CompleteCallback CompleteCallback // 完成重命名所有流程后回调
}

package models

type RenameDst struct {
	Anime    *AnimeEntity
	Content  *TorrentContentItem
	SavePath string
}

type RenameOptions struct {
	Src            string     // 原名
	Dst            *RenameDst // 目标名
	Mode           string
	State          <-chan TorrentState
	RenameCallback func(string) // 重命名完成后回调
	ExitCallback   func()       // 完成重命名所有流程后回调
}

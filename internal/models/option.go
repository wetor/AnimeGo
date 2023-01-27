package models

// =========== Client ===========

type ClientListOptions struct {
	Status   string
	Category string
	Tag      string
}

type ClientAddOptions struct {
	Urls        []string
	SavePath    string
	Category    string
	Tag         string
	SeedingTime int    // 分钟
	Rename      string // 保存名字
}

type ClientDeleteOptions struct {
	Hash       []string
	DeleteFile bool
}

type ClientGetOptions struct {
	Hash string
	Item *TorrentItem
}

type RenameOptions struct {
	Src            string
	Dst            string
	State          <-chan TorrentState
	RenameCallback func() // 重命名完成后回调
	Callback       func() // 完成重命名所有流程后回调
}

// =========== AnimeEntity ===========

type AnimeParseOptions struct {
	Url    string
	Name   string
	Date   string
	Parsed *TitleParsed
}

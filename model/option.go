package model

// =========== Client ===========

type ClientListOptions struct {
	Status   string
	Category string
	Tag      string
}

type ClientRenameOptions struct {
	Hash    string
	OldPath string
	NewPath string
}

type ClientAddOptions struct {
	Urls        []string
	SavePath    string
	Category    string
	Tag         string
	SeedingTime int // 分钟
}

type ClientDeleteOptions struct {
	Hash       []string
	DeleteFile bool
}

type ClientGetOptions struct {
	Hash string
}

// =========== Feed ===========

type FeedParseOptions struct {
	Url          string
	Name         string
	RefreshCache bool // 是否重新下载Url刷新本地缓存
}

// =========== Bangumi ===========

type BangumiParseOptions struct {
	Url  string
	Name string
	ID   int
	Date string
}

// =========== Parser ===========

type ParseBangumiNameOptions struct {
	Name      string
	StartStep int
}

// =========== Process ===========

type ProcessBangumiOptions struct {
	Url  string
	Name string
}

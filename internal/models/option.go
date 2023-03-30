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

// =========== AnimeEntity ===========

type AnimeParseOptions struct {
	Url    string // Mikan url
	Ep     int
	Season int
}

package model

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

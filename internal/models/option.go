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
	Url    string
	Name   string
	Parsed *TitleParsed
}

type PluginFunctionOptions struct {
	Name            string
	ParamsSchema    []string
	ResultSchema    []string
	SkipSchemaCheck bool
}

type PluginVariableOptions struct {
	Name     string
	Nullable bool
}

type PluginLoadOptions struct {
	File      string
	Functions []*PluginFunctionOptions
	Variables []*PluginVariableOptions
}

type PluginExecuteOptions struct {
	File            string
	SkipSchemaCheck bool
}

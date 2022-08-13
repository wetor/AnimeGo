package models

type ParseResult struct {
	NextStep int
	*ParseNameResult
	*ParseEpResult
}
type ParseNameResult struct {
	Name string
}
type ParseEpResult struct {
	Ep int
}

type ParseTagResult struct {
	Resolution string // 分辨率
	Subtitle   string // 字幕语言
	Source     string // 源
}

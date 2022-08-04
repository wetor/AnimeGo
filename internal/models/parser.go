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

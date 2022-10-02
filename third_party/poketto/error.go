package poketto

import "errors"

var (
	CannotParseErr       = errors.New("无法解析")
	CannotParseSeasonErr = errors.New("无法解析季数")
	CannotParseNameErr   = errors.New("无法解析标题")
	CannotParseEpErr     = errors.New("无法解析集数")
	CannotParseTagErr    = errors.New("无法解析标签")
)

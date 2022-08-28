package errors

import "errors"

var (
	ParseAnimeTitleErr = errors.New("解析番剧剧集数失败")
	ParseAnimeNameErr  = errors.New("解析番剧名失败")
)

package errors

import "errors"

var (
	ParseBangumiEpErr   = errors.New("解析番剧剧集数失败")
	ParseBangumiNameErr = errors.New("解析番剧名失败")
)

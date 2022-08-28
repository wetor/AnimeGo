package themoviedb

import "errors"

var (
	NotFoundAnimeNameErr = errors.New("themoviedb: 匹配Seasons失败，番剧名未找到")
	NotMatchSeasonErr    = errors.New("themoviedb: 匹配Seasons失败，可能此番剧未开播")
	ParseAnimeNameErr    = errors.New("themoviedb: 解析番剧名失败")
)

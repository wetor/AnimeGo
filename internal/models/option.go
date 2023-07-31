package models

import "time"

// =========== AnimeEntity ===========

type AnimeParseOverride struct {
	MikanID      int    // Mikan 解析
	BangumiID    int    // 必要，Mikan 解析
	Name         string // 必要，Bangumi 解析
	NameCN       string // 必要，Bangumi 解析
	AirDate      string // 必要，Bangumi 解析，2006-01-02
	Eps          int    // Bangumi 解析
	ThemoviedbID int    // 必要，Themoviedb 解析
	Season       int    // 必要，Themoviedb 解析

}

func (o AnimeParseOverride) OverrideMikan() bool {
	return o.BangumiID != 0
}

func (o AnimeParseOverride) OverrideBangumi() bool {
	if len(o.NameCN) == 0 || len(o.Name) == 0 || len(o.AirDate) == 0 {
		return false
	}
	_, err := time.Parse("2006-01-02", o.AirDate)
	if err != nil {
		return false
	}
	return true
}

func (o AnimeParseOverride) OverrideThemoviedb() bool {
	return o.ThemoviedbID != 0 && o.Season != 0
}

type AnimeParseOptions struct {
	MikanUrl string // Mikan url
	*AnimeParseOverride
}

type ParseOptions struct {
	Title      string
	TorrentUrl string
	MikanUrl   string
	*AnimeParseOverride
}

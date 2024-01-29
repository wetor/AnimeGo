package models

import (
	"sync"
	"time"
)

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
	Input any // Mikan url
	*AnimeParseOverride
}

type ParseOptions struct {
	Title      string // 可选
	TorrentUrl string // 必要
	MikanUrl   string // 和BangumiID二选一
	BangumiID  int    // 和BangumiID二选一，优先
	*AnimeParseOverride
}

type RenamerOptions struct {
	WG            *sync.WaitGroup
	RefreshSecond int
}

type DatabaseOptions struct {
	SavePath string
}

type DownloaderOptions struct {
	RefreshSecond          int
	Category               string
	Tag                    string
	AllowDuplicateDownload bool
	WG                     *sync.WaitGroup
}

type Callback struct {
	Func func(data any) error
}

type NotifierOptions struct {
	DownloadPath string
	SavePath     string
	Rename       string
	Callback     *Callback
}

type FilterOptions struct {
	DelaySecond int
}

type ParserOptions struct {
	TMDBFailSkip           bool
	TMDBFailUseTitleSeason bool
	TMDBFailUseFirstSeason bool
}

type ScheduleOptions struct {
	WG *sync.WaitGroup
}

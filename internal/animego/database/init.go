package database

import (
	"path"
	"sync"

	"github.com/wetor/AnimeGo/pkg/dirdb"
)

const (
	AnimeDBName  = "anime.a_json"
	SeasonDBName = "anime.s_json"
	EpisodeDBFmt = "%s.e_json"
)

var (
	ReInitWG sync.WaitGroup
	Conf     DownloaderConf
)

type DownloaderConf struct {
	RefreshSecond int
	DownloadPath  string
	SavePath      string
	Category      string
	Tag           string
	Rename        string
}

type Options struct {
	DownloaderConf
}

func Init(opts *Options) {
	Conf = opts.DownloaderConf
	dirdb.Init(&dirdb.Options{
		DefaultExt: []string{path.Ext(AnimeDBName), path.Ext(SeasonDBName), path.Ext(EpisodeDBFmt)}, // anime, season
	})
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	opts.DownloaderConf.DownloadPath = Conf.DownloadPath
	opts.DownloaderConf.SavePath = Conf.SavePath

	Conf = opts.DownloaderConf
}

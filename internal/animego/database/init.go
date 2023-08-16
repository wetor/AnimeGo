package database

import "sync"

var CacheMode = true // 将数据库加载到内存中，减少IO

var (
	ReInitWG sync.WaitGroup
	Conf     DownloaderConf
)

type DownloaderConf struct {
	UpdateDelaySecond      int
	DownloadPath           string
	SavePath               string
	Category               string
	Tag                    string
	AllowDuplicateDownload bool
	SeedingTimeMinute      int
	Rename                 string
}

type Options struct {
	DownloaderConf
}

func Init(opts *Options) {
	Conf = opts.DownloaderConf
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	opts.DownloaderConf.DownloadPath = Conf.DownloadPath
	opts.DownloaderConf.SavePath = Conf.SavePath

	Conf = opts.DownloaderConf
}

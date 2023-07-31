package manager

import "sync"

var (
	WG       *sync.WaitGroup
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
	WG *sync.WaitGroup
}

func Init(opts *Options) {
	WG = opts.WG
	Conf = opts.DownloaderConf
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	opts.DownloaderConf.DownloadPath = Conf.DownloadPath
	opts.DownloaderConf.SavePath = Conf.SavePath

	Conf = opts.DownloaderConf
}

package manager

import "sync"

var (
	WG       *sync.WaitGroup
	ReInitWG sync.WaitGroup
	Conf     Downloader
)

type Downloader struct {
	UpdateDelaySecond      int
	DownloadPath           string
	SavePath               string
	Category               string
	Tag                    string
	AllowDuplicateDownload bool
	SeedingTimeMinute      int
	IgnoreSizeMaxKb        int
	Rename                 string
}

type Options struct {
	Downloader
	WG *sync.WaitGroup
}

func Init(opts *Options) {
	WG = opts.WG
	Conf = opts.Downloader
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	opts.Downloader.DownloadPath = Conf.DownloadPath
	opts.Downloader.SavePath = Conf.SavePath

	Conf = opts.Downloader
}

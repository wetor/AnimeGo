package manager

import "sync"

var (
	WG             *sync.WaitGroup
	ReInitWG       sync.WaitGroup
	DownloaderConf Downloader
	FilterConf     Filter
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

type Filter struct {
	MultiGoroutineMax     int
	MultiGoroutineEnabled bool
	UpdateDelayMinute     int
	DelaySecond           int
}

type Options struct {
	Downloader
	Filter
	WG *sync.WaitGroup
}

func Init(opts *Options) {
	WG = opts.WG
	DownloaderConf = opts.Downloader
	FilterConf = opts.Filter
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	opts.Downloader.DownloadPath = DownloaderConf.DownloadPath
	opts.Downloader.SavePath = DownloaderConf.SavePath

	DownloaderConf = opts.Downloader
	FilterConf = opts.Filter
}

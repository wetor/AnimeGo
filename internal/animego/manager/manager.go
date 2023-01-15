package manager

import "sync"

var (
	WG             *sync.WaitGroup
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

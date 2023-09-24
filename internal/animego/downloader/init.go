package downloader

import "sync"

var (
	RefreshSecond          int
	Category               string
	Tag                    string
	AllowDuplicateDownload bool
	SeedingTimeMinute      int
	WG                     *sync.WaitGroup
	ReInitWG               sync.WaitGroup
)

type Options struct {
	RefreshSecond          int
	Category               string
	Tag                    string
	AllowDuplicateDownload bool
	SeedingTimeMinute      int
	WG                     *sync.WaitGroup
}

func Init(opts *Options) {
	RefreshSecond = opts.RefreshSecond
	Category = opts.Category
	Tag = opts.Tag
	AllowDuplicateDownload = opts.AllowDuplicateDownload
	SeedingTimeMinute = opts.SeedingTimeMinute
	WG = opts.WG
}

func ReInit(opts *Options) {
	ReInitWG.Wait()
	RefreshSecond = opts.RefreshSecond
	Category = opts.Category
	Tag = opts.Tag
	AllowDuplicateDownload = opts.AllowDuplicateDownload
	SeedingTimeMinute = opts.SeedingTimeMinute
}
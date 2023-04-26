package torrent

import "os"

var (
	TempPath = os.TempDir()
)

type Options struct {
	TempPath string
}

func Init(opt *Options) {
	TempPath = opt.TempPath
}

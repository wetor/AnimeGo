package torrent

import (
	"os"

	"github.com/wetor/AnimeGo/pkg/xpath"
)

var (
	TempPath = xpath.P(os.TempDir())
)

type Options struct {
	TempPath string
}

func Init(opt *Options) {
	TempPath = opt.TempPath
}

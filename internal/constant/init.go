package constant

import (
	"path"
)

var (
	dataPath = "data"

	CachePath        = "data/cache"
	CacheFile        = "data/cache/bolt.db"
	BangumiCacheFile = "data/cache/bolt_sub.db"
	LogPath          = "data/log"
	LogFile          = "data/log/animego.log"
	PluginPath       = "data/plugin"
	TempPath         = "data/temp"
	WebPath          = "data/web"
)

type Options struct {
	DataPath string
}

func Init(opts *Options) {
	dataPath = opts.DataPath

	CachePath = path.Join(dataPath, "cache")
	CacheFile = path.Join(CachePath, "bolt.db")
	BangumiCacheFile = path.Join(CachePath, "bolt_sub.db")

	LogPath = path.Join(dataPath, "log")
	LogFile = path.Join(LogPath, "animego.log")

	PluginPath = path.Join(dataPath, "plugin")
	TempPath = path.Join(dataPath, "temp")

	WebPath = path.Join(dataPath, "web")
}

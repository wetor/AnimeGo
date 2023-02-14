package constant

import (
	"path/filepath"
)

var (
	AnimeGoGithub = "https://github.com/wetor/AnimeGo"
)

var (
	dataPath         = "data"
	CachePath        = "data/cache"
	CacheFile        = "data/cache/bolt.db"
	BangumiCacheFile = "data/cache/bolt_sub.db"
	LogPath          = "data/log"
	LogFile          = "data/log/animego.log"
	PluginPath       = "data/plugin"
)

type Options struct {
	DataPath string
}

func Init(opts *Options) {
	dataPath = opts.DataPath

	CachePath = filepath.Join(dataPath, "cache")
	CacheFile = filepath.Join(CachePath, "bolt.db")
	BangumiCacheFile = filepath.Join(CachePath, "bolt_sub.db")

	LogPath = filepath.Join(dataPath, "log")
	LogFile = filepath.Join(LogPath, "animego.log")

	PluginPath = filepath.Join(dataPath, "plugin")
}

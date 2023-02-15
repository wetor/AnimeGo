package constant

import "github.com/wetor/AnimeGo/pkg/xpath"

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

	CachePath = xpath.Join(dataPath, "cache")
	CacheFile = xpath.Join(CachePath, "bolt.db")
	BangumiCacheFile = xpath.Join(CachePath, "bolt_sub.db")

	LogPath = xpath.Join(dataPath, "log")
	LogFile = xpath.Join(LogPath, "animego.log")

	PluginPath = xpath.Join(dataPath, "plugin")
}

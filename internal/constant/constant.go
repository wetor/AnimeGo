package constant

import (
	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

var (
	AnimeGoGithub = "https://github.com/wetor/AnimeGo"
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

var (
	PluginTypeBuiltin = assets.BuiltinPrefix
	PluginTypePython  = "python"
)

var (
	PluginTemplatePython   = "python"
	PluginTemplateFilter   = "filter"
	PluginTemplateFeed     = "feed"
	PluginTemplateParser   = "parser"
	PluginTemplateRename   = "rename"
	PluginTemplateSchedule = "schedule"
)

var PluginDirComment = map[string]string{
	PluginTemplateFilter:   "过滤器插件",
	PluginTemplateFeed:     "订阅插件",
	PluginTemplateParser:   "标题解析插件",
	PluginTemplateRename:   "重命名插件",
	PluginTemplateSchedule: "定时任务插件",
}

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
	TempPath = xpath.Join(dataPath, "temp")

	WebPath = xpath.Join(dataPath, "web")
}

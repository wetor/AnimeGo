package constant

import "github.com/wetor/AnimeGo/assets"

const (
	AnimeGoGithub = "https://github.com/wetor/AnimeGo"
)

const (
	DefaultUserAgent = "0.1.0/AnimeGo (https://github.com/wetor/AnimeGo)"
)

const (
	PluginTypeBuiltin = assets.BuiltinPrefix
	PluginTypePython  = "python"
)

const (
	PluginTemplatePython   = "python"
	PluginTemplateFilter   = "filter"
	PluginTemplateFeed     = "feed"
	PluginTemplateRename   = "rename"
	PluginTemplateSchedule = "schedule"
)

const (
	WriteFilePerm = 0664 // rw,rw,r
)

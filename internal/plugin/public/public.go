package public

import (
	"path/filepath"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/utils"
)

var (
	py         = &python.Python{}
	pluginPath string
)

type Options struct {
	PluginPath string
}

func Init(opt *Options) {
	pluginPath = opt.PluginPath
}

func ParserName(title string) (ep *models.TitleParsed) {
	pluginFile := filepath.Join(pluginPath, "lib/Auto_Bangumi/raw_parser.py")
	if !utils.IsExist(pluginFile) {
		utils.CopyDir(assets.Plugin, "plugin", pluginPath, true, false)
	}
	py.Load(&models.PluginLoadOptions{
		File: pluginFile,
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "parse",
				SkipSchemaCheck: true,
			},
		},
	})
	result := py.Run("parse", models.Object{
		"title": title,
	})
	ep = &models.TitleParsed{
		TitleRaw: title,
	}
	utils.Map2ModelByJson(result, ep)
	if len(ep.NameCN) > 0 {
		ep.Name = ep.NameCN
	} else if len(ep.Name) == 0 && len(ep.NameEN) > 0 {
		ep.Name = ep.NameEN
	}
	return ep
}

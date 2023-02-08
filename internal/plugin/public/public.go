package public

import (
	"path"

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
	pluginFile := path.Join(pluginPath, "lib/Auto_Bangumi/raw_parser.py")
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

	if tmp, has := result["title_en"]; has {
		ep.NameEN = tmp.(string)
		ep.Name = ep.NameEN
	}
	if tmp, has := result["title_jp"]; has {
		ep.Name = tmp.(string)
	}
	if tmp, has := result["title_zh"]; has {
		ep.NameCN = tmp.(string)
		ep.Name = ep.NameCN
	}
	if tmp, has := result["season"]; has {
		ep.Season = int(tmp.(int64))
	}
	if tmp, has := result["season_raw"]; has {
		ep.SeasonRaw = tmp.(string)
	}
	if tmp, has := result["episode"]; has {
		ep.Ep = int(tmp.(int64))
	}
	if tmp, has := result["sub"]; has {
		ep.Sub = tmp.(string)
	}
	if tmp, has := result["group"]; has {
		ep.Group = tmp.(string)
	}
	if tmp, has := result["resolution"]; has {
		ep.Definition = tmp.(string)
	}
	if tmp, has := result["source"]; has {
		ep.Source = tmp.(string)
	}
	return ep
}

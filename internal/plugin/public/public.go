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
	pluginFile := path.Join(pluginPath, "anisource/Auto_Bangumi/raw_parser.py")
	if !utils.IsExist(pluginFile) {
		utils.CopyDir(assets.Plugin, "plugin", pluginPath, true, true)
	}
	py.SetSchema([]string{"title"}, []string{})
	result := py.Execute(pluginFile, models.Object{
		"title": title,
	})
	ep = &models.TitleParsed{
		TitleRaw: title,
	}
	if obj, ok := result.(models.Object); ok {
		if tmp, has := obj["title_en"]; has {
			ep.NameEN = tmp.(string)
			ep.Name = ep.NameEN
		}
		if tmp, has := obj["title_jp"]; has {
			ep.Name = tmp.(string)
		}
		if tmp, has := obj["title_zh"]; has {
			ep.NameCN = tmp.(string)
			ep.Name = ep.NameCN
		}
		if tmp, has := obj["season"]; has {
			ep.Season = int(tmp.(int64))
		}
		if tmp, has := obj["season_raw"]; has {
			ep.SeasonRaw = tmp.(string)
		}
		if tmp, has := obj["episode"]; has {
			ep.Ep = int(tmp.(int64))
		}
		if tmp, has := obj["sub"]; has {
			ep.Sub = tmp.(string)
		}
		if tmp, has := obj["group"]; has {
			ep.Group = tmp.(string)
		}
		if tmp, has := obj["resolution"]; has {
			ep.Definition = tmp.(string)
		}
		if tmp, has := obj["source"]; has {
			ep.Source = tmp.(string)
		}
	}
	return ep
}

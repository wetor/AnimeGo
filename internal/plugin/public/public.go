package public

import (
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	parserScript = "filter/Auto_Bangumi/raw_parser.py"
	FuncMain     = "main"
)

var (
	py api.Plugin = nil
)

func ParserName(title string) (ep *models.TitleParsed) {
	if py == nil {
		py = plugin.LoadPlugin(&plugin.LoadPluginOptions{
			Plugin: &models.Plugin{
				Type: "builtin",
				File: "raw_parser.py",
			},
			EntryFunc: FuncMain,
			FuncSchema: []*pkgPlugin.FuncSchemaOptions{
				{
					Name:            FuncMain,
					SkipSchemaCheck: true,
				},
			},
		})
	}
	result := py.Run(FuncMain, map[string]any{
		"title": title,
	})
	ep = &models.TitleParsed{
		TitleRaw: title,
	}
	utils.MapToStruct(result, ep)
	if len(ep.NameCN) > 0 {
		ep.Name = ep.NameCN
	} else if len(ep.Name) == 0 && len(ep.NameEN) > 0 {
		ep.Name = ep.NameEN
	}
	return ep
}

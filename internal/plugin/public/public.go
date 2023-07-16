package public

import (
	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	FuncMain = "main"
)

var (
	py api.Plugin = nil
)

func ParserName(title string) (ep *models.TitleParsed, err error) {
	if py == nil {
		py, err = plugin.LoadPlugin(&plugin.LoadPluginOptions{
			Plugin: &models.Plugin{
				Type: "builtin",
				File: assets.BuiltinRawParser,
			},
			EntryFunc: FuncMain,
			FuncSchema: []*pkgPlugin.FuncSchemaOptions{
				{
					Name:            FuncMain,
					SkipSchemaCheck: true,
				},
			},
		})
		if err != nil {
			return nil, err
		}
	}
	result, err := py.Run(FuncMain, map[string]any{
		"title": title,
	})
	if err != nil {
		return nil, err
	}
	ep = &models.TitleParsed{
		TitleRaw: title,
	}
	err = utils.MapToStruct(result, ep)
	if err != nil {
		return nil, err
	}
	if len(ep.NameCN) > 0 {
		ep.Name = ep.NameCN
	} else if len(ep.Name) == 0 && len(ep.NameEN) > 0 {
		ep.Name = ep.NameEN
	}
	return ep, nil
}

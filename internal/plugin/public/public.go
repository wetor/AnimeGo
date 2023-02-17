package public

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin/python"
	"github.com/wetor/AnimeGo/internal/utils"
)

const (
	parserScript = "filter/Auto_Bangumi/raw_parser.py"
)

var (
	py = &python.Python{}
)

func ParserName(title string) (ep *models.TitleParsed) {
	py.Load(&models.PluginLoadOptions{
		File: parserScript,
		Functions: []*models.PluginFunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})
	result := py.Run("main", models.Object{
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

package public

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/plugin/python"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	parserScript = "filter/Auto_Bangumi/raw_parser.py"
)

var (
	py = &python.Python{}
)

func ParserName(title string) (ep *models.TitleParsed) {
	py.Load(&plugin.LoadOptions{
		File: parserScript,
		Functions: []*plugin.FunctionOptions{
			{
				Name:            "main",
				SkipSchemaCheck: true,
			},
		},
	})
	result := py.Run("main", map[string]any{
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

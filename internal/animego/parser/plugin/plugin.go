package plugin

import (
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	FuncParse = "parse"
)

type Parser struct {
	plugin         *models.Plugin
	pluginInstance api.Plugin
	single         bool
}

func NewParserPlugin(pluginInfo *models.Plugin, single bool) *Parser {
	return &Parser{
		plugin: pluginInfo,
		single: single,
	}
}

func (p *Parser) Parse(title string) *models.TitleParsed {
	if p.pluginInstance == nil || !p.single {
		p.pluginInstance = plugin.LoadPlugin(&plugin.LoadPluginOptions{
			Plugin:    p.plugin,
			EntryFunc: FuncParse,
			FuncSchema: []*pkgPlugin.FuncSchemaOptions{
				{
					Name:         FuncParse,
					ParamsSchema: []string{"title"},
					ResultSchema: []string{"error", "data"},
					DefaultArgs:  p.plugin.Args,
				},
			},
		})
	}
	result := p.pluginInstance.Run(FuncParse, map[string]any{
		"title": title,
	})
	if result["error"] != nil {
		log.Debugf("", errors.NewAniErrorD(result["error"]))
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", p.plugin.Type, p.plugin.File, result["error"])
	}
	ep := &models.TitleParsed{
		TitleRaw: title,
	}
	data, ok := result["data"].(map[string]any)
	if !ok {
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", p.plugin.Type, p.plugin.File, result["data"])
		return ep
	}
	utils.MapToStruct(data, ep)
	return ep
}

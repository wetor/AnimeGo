package parser

import (
	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
)

var PluginSet = wire.NewSet(
	NewParserPlugin,
)

type Parser struct {
	plugin         *models.Plugin
	pluginInstance api.Plugin
	single         bool
}

func NewParserPlugin(pluginInfo *models.Plugin) *Parser {
	return &Parser{
		plugin: pluginInfo,
		single: true,
	}
}

func (p *Parser) Parse(title string) (*models.TitleParsed, error) {
	var err error
	if p.pluginInstance == nil || !p.single {
		p.pluginInstance, err = plugin.LoadPlugin(&plugin.LoadPluginOptions{
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
		if err != nil {
			return nil, err
		}
	}
	result, err := p.pluginInstance.Run(FuncParse, map[string]any{
		"title": title,
	})
	if err != nil {
		return nil, err
	}
	if result["error"] != nil {
		err = errors.WithStack(&exceptions.ErrPlugin{Type: p.plugin.Type, File: p.plugin.File, Message: result["error"]})
		log.DebugErr(err)
		log.Warnf("[Plugin] %s插件(%s)执行错误: %v", p.plugin.Type, p.plugin.File, result["error"])
		return nil, err
	}
	ep := &models.TitleParsed{
		TitleRaw: title,
	}
	data := result["data"].(map[string]any)
	err = utils.MapToStruct(data, ep)
	if err != nil {
		log.DebugErr(err)
		return nil, errors.WithStack(&exceptions.ErrPlugin{Type: p.plugin.Type, File: p.plugin.File, Message: "类型转换错误"})
	}
	return ep, nil
}

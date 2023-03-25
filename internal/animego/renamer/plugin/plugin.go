package plugin

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
)

const FuncRename = "rename"

type Rename struct {
	plugin *models.Plugin
}

func NewRenamePlugin(pluginInfo []models.Plugin) *Rename {
	for _, p := range pluginInfo {
		if p.Enable {
			return &Rename{
				plugin: &p,
			}
		}
	}
	return &Rename{}
}

func (p *Rename) Rename(anime *models.AnimeEntity, src string) string {
	pluginInstance := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    p.plugin,
		EntryFunc: FuncRename,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         FuncRename,
				ParamsSchema: []string{"anime", "src"},
				ResultSchema: []string{"error", "dst"},
				DefaultArgs:  p.plugin.Args,
			},
		},
	})
	result := pluginInstance.Run(FuncRename, map[string]any{
		"anime": anime,
		"src":   src,
	})
	if result["error"] != nil {
		log.Debugf("", errors.NewAniErrorD(result["error"]))
		log.Warnf("[Plugin] Rename插件(%s)执行错误: %v", p.plugin.File, result["error"])
	}
	if dst, ok := result["dst"].(string); ok && len(dst) != 0 {
		log.Debugf("[Plugin] Rename插件(%s): %s -> %s", p.plugin.File, src, dst)
		return dst
	}
	return ""
}

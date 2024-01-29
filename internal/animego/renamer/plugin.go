package renamer

import (
	"path"

	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

var PluginSet = wire.NewSet(
	NewRenamePlugin,
)

type Rename struct {
	plugin *models.Plugin
}

func NewRenamePlugin(pluginInfo *models.Plugin) *Rename {
	return &Rename{
		plugin: pluginInfo,
	}
}

func (p *Rename) Rename(anime *models.AnimeEntity, index int, filename string) (*models.RenameResult, error) {
	pluginInstance, err := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    p.plugin,
		EntryFunc: FuncRename,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         FuncRename,
				ParamsSchema: []string{"anime", "filename"},
				ResultSchema: []string{"error", "filename", "dir,optional"},
				DefaultArgs:  p.plugin.Args,
			},
		},
		VarSchema: []*pkgPlugin.VarSchemaOptions{
			{
				Name:     VarScrape,
				Nullable: true,
			},
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "加载重命名插件失败")
	}

	obj := utils.StructToMap(anime)
	obj["ep_type"] = anime.Ep[index].Type
	obj["ep"] = anime.Ep[index].Ep
	result, err := pluginInstance.Run(FuncRename, map[string]any{
		"anime":    obj,
		"filename": filename,
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

	val, err := pluginInstance.Get(VarScrape)
	if err != nil {
		return nil, err
	}
	scrape := true
	if val != nil {
		scrape, _ = val.(bool)
	}

	renameResult := &models.RenameResult{
		Index:  index,
		Scrape: scrape,
	}
	if dst, ok := result["filename"].(string); ok && len(dst) != 0 {
		renameResult.Filename = dst
		log.Debugf("[Plugin] Rename插件(%s): %s -> %s", p.plugin.File, filename, dst)
	}
	if dir, ok := result["dir"].(string); ok {
		renameResult.AnimeDir = dir
	} else {
		renameResult.AnimeDir = xpath.Root(renameResult.Filename)
		if renameResult.AnimeDir == "." {
			renameResult.AnimeDir = ""
		}
	}
	renameResult.SeasonDir = path.Dir(renameResult.Filename)
	return renameResult, nil
}

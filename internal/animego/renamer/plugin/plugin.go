package plugin

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	pkgPlugin "github.com/wetor/AnimeGo/pkg/plugin"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	FuncRename     = "rename"
	VarWriteTVShow = "__write_tvshow__"
)

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

func (p *Rename) Rename(anime *models.AnimeEntity, filename string) *models.RenameResult {
	pluginInstance := plugin.LoadPlugin(&plugin.LoadPluginOptions{
		Plugin:    p.plugin,
		EntryFunc: FuncRename,
		FuncSchema: []*pkgPlugin.FuncSchemaOptions{
			{
				Name:         FuncRename,
				ParamsSchema: []string{"anime", "filename", "default_filepath"},
				ResultSchema: []string{"error", "filepath", "tvshow_dir,optional"},
				DefaultArgs:  p.plugin.Args,
			},
		},
		VarSchema: []*pkgPlugin.VarSchemaOptions{
			{
				Name:     VarWriteTVShow,
				Nullable: true,
			},
		},
	})
	result := pluginInstance.Run(FuncRename, map[string]any{
		"anime":            anime,
		"filename":         filename,
		"default_filepath": xpath.Join(anime.DirName(), anime.FileName()+xpath.Ext(filename)),
	})
	if result["error"] != nil {
		log.Debugf("", errors.NewAniErrorD(result["error"]))
		log.Warnf("[Plugin] Rename插件(%s)执行错误: %v", p.plugin.File, result["error"])
		return nil
	}

	tvshow := true
	val := pluginInstance.Get(VarWriteTVShow)
	if val != nil {
		tvshow = val.(bool)
	}

	renameResult := &models.RenameResult{}
	if dst, ok := result["filepath"].(string); ok && len(dst) != 0 {
		renameResult.Filepath = dst
		log.Debugf("[Plugin] Rename插件(%s): %s -> %s", p.plugin.File, filename, dst)
	}
	if tvshow {
		if dir, ok := result["tvshow_dir"].(string); ok {
			renameResult.TVShowDir = dir
		} else {
			renameResult.TVShowDir = xpath.Join(renameResult.Filepath, "../..")
			if renameResult.TVShowDir == "." {
				renameResult.TVShowDir = ""
			}
		}
	}
	return renameResult
}

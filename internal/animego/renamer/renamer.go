package renamer

import (
	"github.com/wetor/AnimeGo/internal/api"
	"os"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

type Renamer struct {
	plugin api.RenamerPlugin
}

func NewRenamer(plugin api.RenamerPlugin) *Renamer {
	return &Renamer{
		plugin: plugin,
	}
}

func (r *Renamer) stateSeeding(mode string, src, dst, renamePath string, callback func(string)) bool {
	var err error
	switch mode {
	case "link_delete":
		log.Infof("[重命名] 链接「%s」->「%s」", src, dst)
		err = utils.CreateLink(src, dst)
		errors.NewAniErrorD(err).TryPanic()
		callback(renamePath)
	case "link":
		log.Infof("[重命名] 链接「%s」->「%s」", src, dst)
		err = utils.CreateLink(src, dst)
		errors.NewAniErrorD(err).TryPanic()
		callback(renamePath)
		return true
	case "move":
		log.Infof("[重命名] 移动「%s」->「%s」", src, dst)
		err = utils.Rename(src, dst)
		errors.NewAniErrorD(err).TryPanic()
		callback(renamePath)
		return true
	}
	return false
}

func (r *Renamer) stateComplete(mode string, src, dst, renamePath string, callback func(string)) bool {
	var err error
	switch mode {
	case "wait_move":
		log.Infof("[重命名] 移动「%s」->「%s」", src, dst)
		err = utils.Rename(src, dst)
		errors.NewAniErrorD(err).TryPanic()
		callback(renamePath)
		return true
	case "link_delete":
		if !utils.IsExist(dst) {
			log.Infof("[重命名] 链接「%s」->「%s」", src, dst)
			err = utils.CreateLink(src, dst)
			errors.NewAniErrorD(err).TryPanic()
			callback(renamePath)
		}
		log.Infof("[重命名] 删除「%s」", src)
		err = os.Remove(src)
		errors.NewAniErrorD(err).TryPanic()
		return true
	}
	return false
}

func (r *Renamer) Rename(opt *models.RenameOptions) {
	var renamePath string
	if r.plugin != nil {
		renamePath = r.plugin.Rename(opt.Dst.Anime, opt.Dst.Content.Name)
	}
	if len(renamePath) == 0 {
		renamePath = xpath.Join(opt.Dst.Anime.DirName(), opt.Dst.Anime.FileName()+xpath.Ext(opt.Dst.Content.Name))
	}
	dst := xpath.Join(opt.Dst.SavePath, renamePath)
	go func() {
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
				})
				state := <-opt.State

				existSrc := utils.IsExist(opt.Src)
				existDst := utils.IsExist(dst)

				if !existSrc && !existDst {
					errors.NewAniError("未找到文件：" + opt.Src).TryPanic()
				} else if !existSrc && existDst {
					// 已经移动完成
					log.Warnf("[跳过重命名] 可能已经移动完成「%s」->「%s」", opt.Src, dst)
					if state == downloader.StateComplete {
						exit = r.stateComplete(opt.Mode, opt.Src, dst, renamePath, opt.RenameCallback)
						return
					}
				}
				if state == downloader.StateSeeding {
					exit = r.stateSeeding(opt.Mode, opt.Src, dst, renamePath, opt.RenameCallback)
				} else if state == downloader.StateComplete {
					exit = r.stateSeeding(opt.Mode, opt.Src, dst, renamePath, opt.RenameCallback)
					exit = r.stateComplete(opt.Mode, opt.Src, dst, renamePath, opt.RenameCallback)
				}
			}()
			if exit {
				opt.ExitCallback()
				return
			}
		}
	}()
}

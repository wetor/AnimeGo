package downloader

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"os"
)

func RenameAnime(opt *models.RenameOptions) {
	var err error
	rename := store.Config.Advanced.Download.Rename
	go func() {
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					zap.S().Error(err)
				})
				state := <-opt.State

				if !utils.IsExist(opt.Src) {
					errors.NewAniError("未找到文件：" + opt.Src).TryPanic()
				}

				switch rename {
				case "link", "link_delete":
					zap.S().Infof("[重命名] 链接「%s」->「%s」", opt.Src, opt.Dst)
					err = utils.CreateLink(opt.Src, opt.Dst)
					errors.NewAniErrorD(err).TryPanic()
					opt.Callback()
					if rename == "link" {
						exit = true
					}
				case "move":
					zap.S().Infof("[重命名] 移动「%s」->「%s」", opt.Src, opt.Dst)
					err = utils.Rename(opt.Src, opt.Dst)
					errors.NewAniErrorD(err).TryPanic()
					opt.Callback()
					exit = true
				}

				if state == StateComplete {
					switch rename {
					case "wait_move":
						zap.S().Infof("[重命名] 移动「%s」->「%s」", opt.Src, opt.Dst)
						err = utils.Rename(opt.Src, opt.Dst)
						errors.NewAniErrorD(err).TryPanic()
						opt.Callback()
						exit = true
					case "link_delete":
						if !utils.IsExist(opt.Dst) {
							zap.S().Infof("[重命名] 链接「%s」->「%s」", opt.Src, opt.Dst)
							err = utils.CreateLink(opt.Src, opt.Dst)
							errors.NewAniErrorD(err).TryPanic()
							opt.Callback()
						}
						zap.S().Infof("[重命名] 删除「%s」", opt.Src)
						err = os.Remove(opt.Src)
						errors.NewAniErrorD(err).TryPanic()
						exit = true
					}
				}
			}()
			if exit {
				return
			}
		}
	}()
}

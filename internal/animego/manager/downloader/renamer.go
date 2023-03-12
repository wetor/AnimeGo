package downloader

import (
	"os"

	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
)

type RenameOptions struct {
	Src            string
	Dst            string
	State          <-chan models.TorrentState
	RenameCallback func() // 重命名完成后回调
	Callback       func() // 完成重命名所有流程后回调
}

func RenameAnime(opt *RenameOptions) {
	if opt.Src == opt.Dst {
		return
	}
	var err error
	rename := manager.DownloaderConf.Rename
	go func() {
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
				})
				state := <-opt.State

				existSrc := utils.IsExist(opt.Src)
				existDst := utils.IsExist(opt.Dst)

				if !existSrc && !existDst {
					errors.NewAniError("未找到文件：" + opt.Src).TryPanic()
				} else if !existSrc && existDst {
					// 已经移动完成
					log.Warnf("[跳过重命名] 可能已经移动完成「%s」->「%s」", opt.Src, opt.Dst)
					opt.RenameCallback()
					if state == StateComplete {
						exit = true
						return
					}
				}

				switch rename {
				case "link", "link_delete":
					log.Infof("[重命名] 链接「%s」->「%s」", opt.Src, opt.Dst)
					err = utils.CreateLink(opt.Src, opt.Dst)
					errors.NewAniErrorD(err).TryPanic()
					opt.RenameCallback()
					if rename == "link" {
						exit = true
					}
				case "move":
					log.Infof("[重命名] 移动「%s」->「%s」", opt.Src, opt.Dst)
					err = utils.Rename(opt.Src, opt.Dst)
					errors.NewAniErrorD(err).TryPanic()
					opt.RenameCallback()
					exit = true
				}

				if state == StateComplete {
					switch rename {
					case "wait_move":
						log.Infof("[重命名] 移动「%s」->「%s」", opt.Src, opt.Dst)
						err = utils.Rename(opt.Src, opt.Dst)
						errors.NewAniErrorD(err).TryPanic()
						opt.RenameCallback()
						exit = true
					case "link_delete":
						if !utils.IsExist(opt.Dst) {
							log.Infof("[重命名] 链接「%s」->「%s」", opt.Src, opt.Dst)
							err = utils.CreateLink(opt.Src, opt.Dst)
							errors.NewAniErrorD(err).TryPanic()
							opt.RenameCallback()
						}
						log.Infof("[重命名] 删除「%s」", opt.Src)
						err = os.Remove(opt.Src)
						errors.NewAniErrorD(err).TryPanic()
						exit = true
					}
				}
			}()
			if exit {
				opt.Callback()
				return
			}
		}
	}()
}

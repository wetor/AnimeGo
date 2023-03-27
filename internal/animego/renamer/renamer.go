package renamer

import (
	"context"
	"os"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

type RenameTask struct {
	// 只读
	Src              string            // 原名
	RenameDst        *models.RenameDst // 目标名
	Mode             string
	State            <-chan models.TorrentState
	RenameCallback   models.RenameCallback   // 重命名完成后回调
	CompleteCallback models.CompleteCallback // 完成重命名所有流程后回调
	Dst              string
	RenameResult     *models.RenameResult

	// 读写
	Complete bool // 是否完成任务
}

type Renamer struct {
	plugin api.RenamerPlugin
	tasks  map[string]*RenameTask
}

func NewRenamer(plugin api.RenamerPlugin) *Renamer {
	return &Renamer{
		plugin: plugin,
		tasks:  make(map[string]*RenameTask),
	}
}

func (r *Renamer) stateSeeding(task *RenameTask) (complete bool) {
	var err error
	switch task.Mode {
	case "wait_move":
		complete = false
	case "link_delete":
		log.Infof("[重命名] 链接「%s」->「%s」", task.Src, task.Dst)
		err = utils.CreateLink(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = false
		task.RenameCallback(task.RenameResult)
	case "link":
		log.Infof("[重命名] 链接「%s」->「%s」", task.Src, task.Dst)
		err = utils.CreateLink(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
		task.RenameCallback(task.RenameResult)
	case "move":
		log.Infof("[重命名] 移动「%s」->「%s」", task.Src, task.Dst)
		err = utils.Rename(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
		task.RenameCallback(task.RenameResult)
	default:
		errors.NewAniErrorf("不支持的重命名模式 %s", task.Mode).TryPanic()
	}
	return
}

func (r *Renamer) stateComplete(task *RenameTask) (complete bool) {
	var err error
	switch task.Mode {
	case "wait_move":
		log.Infof("[重命名] 移动「%s」->「%s」", task.Src, task.Dst)
		err = utils.Rename(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
		task.RenameCallback(task.RenameResult)
	case "link_delete":
		if !utils.IsExist(task.Dst) {
			r.stateSeeding(task)
		}
		log.Infof("[重命名] 删除「%s」", task.Src)
		err = os.Remove(task.Src)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
	case "link":
		complete = true
	case "move":
		complete = true
	default:
		errors.NewAniErrorf("不支持的重命名模式 %s", task.Mode).TryPanic()
	}
	return
}

func (r *Renamer) AddRenameTask(opt *models.RenameOptions) {
	var result *models.RenameResult
	if r.plugin != nil {
		result = r.plugin.Rename(opt.Dst.Anime, opt.Dst.Content.Name)
	}
	if result == nil {
		result = &models.RenameResult{
			Filepath:  xpath.Join(opt.Dst.Anime.DirName(), opt.Dst.Anime.FileName()+xpath.Ext(opt.Dst.Content.Name)),
			TVShowDir: opt.Dst.Anime.DirName(),
		}
	}
	dst := xpath.Join(opt.Dst.SavePath, result.Filepath)
	r.tasks[opt.Src] = &RenameTask{
		Src:              opt.Src,
		RenameDst:        opt.Dst,
		Mode:             opt.Mode,
		State:            opt.State,
		RenameCallback:   opt.RenameCallback,
		CompleteCallback: opt.CompleteCallback,
		Dst:              dst,
		RenameResult:     result,
		Complete:         false,
	}
}

func (r *Renamer) Update(ctx context.Context) {
	var deleteKeys []string
	for name, task := range r.tasks {
		if task.Complete {
			deleteKeys = append(deleteKeys, name)
			continue
		}
		select {
		case state := <-task.State:
			go func(task *RenameTask) {
				defer func() {
					if task.Complete {
						task.CompleteCallback()
					}
				}()
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
				})
				existSrc := utils.IsExist(task.Src)
				existDst := utils.IsExist(task.Dst)

				if !existSrc && !existDst {
					errors.NewAniError("未找到文件：" + task.Src).TryPanic()
				} else if !existSrc && existDst {
					// 已经移动完成
					log.Warnf("[跳过重命名] 可能已经移动完成「%s」->「%s」", task.Src, task.Dst)
					if state == downloader.StateComplete {
						task.Complete = r.stateComplete(task)
						return
					}
				}
				if state == downloader.StateSeeding {
					task.Complete = r.stateSeeding(task)
				} else if state == downloader.StateComplete {
					task.Complete = r.stateSeeding(task)
					task.Complete = r.stateComplete(task)
				}
			}(task)
		default:

		}
	}

	for _, k := range deleteKeys {
		delete(r.tasks, k)
	}
}

func (r *Renamer) sleep(ctx context.Context) {
	utils.Sleep(UpdateDelaySecond, ctx)
}

func (r *Renamer) Start(ctx context.Context) {
	WG.Add(1)
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
					r.sleep(ctx)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 renamer")
					exit = true
					return
				default:
					r.Update(ctx)
					r.sleep(ctx)
				}
			}()
			if exit {
				return
			}
		}
	}()
}

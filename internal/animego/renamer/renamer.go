package renamer

import (
	"context"
	"os"
	"sync"

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
	Src            string // 原名
	Dst            string
	Mode           string
	State          <-chan models.TorrentState
	RenameCallback models.RenameCallback // 重命名完成后回调
	RenameResult   *models.RenameResult

	// 读写
	Complete bool // 是否完成任务
}

type RenameTaskGroup struct {
	Tasks            []*RenameTask
	CompleteCallback models.CompleteCallback // 完成重命名所有流程后回调
}

func (r *RenameTaskGroup) Complete() bool {
	for _, t := range r.Tasks {
		if !t.Complete {
			return false
		}
	}
	return true
}

type Manager struct {
	plugin     api.RenamerPlugin
	taskGroups map[string]*RenameTaskGroup
	sync.Mutex
}

func NewManager(plugin api.RenamerPlugin) *Manager {
	return &Manager{
		plugin:     plugin,
		taskGroups: make(map[string]*RenameTaskGroup),
	}
}

func (m *Manager) stateSeeding(task *RenameTask) (complete bool) {
	var err error
	switch task.Mode {
	case "wait_move":
		complete = false
	case "link_delete":
		log.Infof("[重命名] 链接「%s」->「%s」", task.Src, task.Dst)
		err = utils.CreateLink(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = false
	case "link":
		log.Infof("[重命名] 链接「%s」->「%s」", task.Src, task.Dst)
		err = utils.CreateLink(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
	case "move":
		log.Infof("[重命名] 移动「%s」->「%s」", task.Src, task.Dst)
		err = utils.Rename(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
	default:
		errors.NewAniErrorf("不支持的重命名模式 %s", task.Mode).TryPanic()
	}
	return
}

func (m *Manager) stateComplete(task *RenameTask) (complete bool) {
	var err error
	switch task.Mode {
	case "wait_move":
		log.Infof("[重命名] 移动「%s」->「%s」", task.Src, task.Dst)
		err = utils.Rename(task.Src, task.Dst)
		errors.NewAniErrorD(err).TryPanic()
		complete = true
	case "link_delete":
		if !utils.IsExist(task.Dst) {
			m.stateSeeding(task)
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

func (m *Manager) HasRenameTask(name string) bool {
	_, ok := m.taskGroups[name]
	return ok
}

func (m *Manager) AddRenameTask(opt *models.RenameOptions) {
	m.Lock()
	defer m.Unlock()
	name := opt.Entity.FullName()
	srcFiles := opt.Entity.FilePathSrc()
	dstFiles := opt.Entity.FilePath()
	m.taskGroups[name] = &RenameTaskGroup{
		Tasks:            make([]*RenameTask, len(opt.Entity.Ep)),
		CompleteCallback: opt.CompleteCallback,
	}
	for i := range opt.Entity.Ep {
		var result *models.RenameResult
		if m.plugin != nil {
			result = m.plugin.Rename(opt.Entity, i, xpath.Base(srcFiles[i]))
		}
		if result == nil {
			result = &models.RenameResult{
				Filepath:  dstFiles[i],
				TVShowDir: opt.Entity.DirName(),
			}
		}
		src := xpath.Join(opt.SrcDir, srcFiles[i])
		dst := xpath.Join(opt.DstDir, result.Filepath)
		result.Index = i
		m.taskGroups[name].Tasks[i] = &RenameTask{
			Src:            src,
			Dst:            dst,
			Mode:           opt.Mode,
			State:          opt.State[i],
			RenameCallback: opt.RenameCallback,
			RenameResult:   result,
			Complete:       false,
		}
	}
}

func (m *Manager) Update(ctx context.Context) {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	var deleteKeys []string

	for name, taskGroup := range m.taskGroups {
		if taskGroup.Complete() {
			taskGroup.CompleteCallback(nil)
			deleteKeys = append(deleteKeys, name)
			continue
		}
		for _, task := range taskGroup.Tasks {
			if task.Complete {
				continue
			}
			select {
			case state := <-task.State:
				go func(task *RenameTask) {
					defer func() {
						if task.Complete {
							task.RenameCallback(task.RenameResult)
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
							task.Complete = m.stateComplete(task)
							return
						}
					}
					if state == downloader.StateSeeding {
						task.Complete = m.stateSeeding(task)
					} else if state == downloader.StateComplete {
						task.Complete = m.stateSeeding(task)
						task.Complete = m.stateComplete(task)
					}
				}(task)
			default:
			}
		}
	}

	for _, k := range deleteKeys {
		delete(m.taskGroups, k)
	}
}

func (m *Manager) sleep(ctx context.Context) {
	utils.Sleep(UpdateDelaySecond, ctx)
}

func (m *Manager) Start(ctx context.Context) {
	WG.Add(1)
	// 刷新信息、接收下载、接收退出指令协程
	go func() {
		defer WG.Done()
		for {
			exit := false
			func() {
				defer errors.HandleError(func(err error) {
					log.Errorf("", err)
					m.sleep(ctx)
				})
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 renamer")
					exit = true
					return
				default:
					m.Update(ctx)
					m.sleep(ctx)
				}
			}()
			if exit {
				return
			}
		}
	}()
}

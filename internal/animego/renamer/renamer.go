package renamer

import (
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	RenameStateError = iota - 1
	RenameStateStart
	RenameStateSeeding
	RenameStateComplete
	RenameStateEnd
)

const (
	AllRenameStateError = iota - 1
	AllRenameStateStart
	AllRenameStateIncomplete
	AllRenameStateComplete
)

const (
	RenameStateChanCap = 5
	RenameMaxErrCount  = 3
)

type RenameTask struct {
	// 只读
	Src            string // 原名
	Dst            string
	Mode           string
	StateChan      chan models.TorrentState
	RenameCallback models.RenameCallback // 重命名完成后回调
	RenameResult   *models.RenameResult

	// 读写
	RenameState int
	State       models.TorrentState
	ErrCount    int
}

type RenameTaskGroup struct {
	Tasks            []*RenameTask
	CompleteCallback models.CompleteCallback // 完成重命名所有流程后回调
}

func (r *RenameTaskGroup) Complete() int {
	incomplete := 0
	for _, t := range r.Tasks {
		if t.RenameState != RenameStateEnd {
			incomplete++
		}
	}
	if incomplete == 0 {
		return AllRenameStateComplete
	} else if incomplete == len(r.Tasks) {
		return AllRenameStateStart
	} else {
		return AllRenameStateIncomplete
	}
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

func (m *Manager) SetDownloadState(name string, epIndex int, state models.TorrentState) error {
	g, ok := m.taskGroups[name]
	if !ok {
		return errors.WithStack(exceptions.ErrRename{Src: name, Message: "任务不存在"})
	}
	if epIndex < 0 || epIndex >= len(g.Tasks) {
		return errors.WithStack(exceptions.ErrRename{Src: fmt.Sprintf("%s[%d]", name, epIndex), Message: "任务不存在"})
	}
	g.Tasks[epIndex].StateChan <- state
	return nil
}

func (m *Manager) HasRenameTask(name string) bool {
	_, ok := m.taskGroups[name]
	return ok
}

func (m *Manager) GetRenameTaskState(name string) (int, error) {
	g, ok := m.taskGroups[name]
	if !ok {
		return AllRenameStateError, errors.WithStack(exceptions.ErrRename{Src: name, Message: "任务不存在"})
	}
	return g.Complete(), nil
}

func (m *Manager) GetEpTaskState(name string, epIndex int) (int, error) {
	g, ok := m.taskGroups[name]
	if !ok {
		return RenameStateError, errors.WithStack(exceptions.ErrRename{Src: name, Message: "任务不存在"})
	}
	if epIndex < 0 || epIndex >= len(g.Tasks) {
		return RenameStateError, errors.WithStack(exceptions.ErrRename{Src: fmt.Sprintf("%s[%d]", name, epIndex), Message: "任务不存在"})
	}
	return g.Tasks[epIndex].RenameState, nil
}

func (m *Manager) stateSeeding(task *RenameTask) (err error) {
	if task.ErrCount >= RenameMaxErrCount {
		log.Warnf("[重命名] 失败，跳过流程：「%s」->「%s」", task.Src, task.Dst)
		task.RenameState = RenameStateEnd
		return nil
	}
	defer func() {
		if err != nil {
			task.ErrCount++
		}
	}()
	switch task.Mode {
	case "wait_move":
		task.RenameState = RenameStateComplete
	case "link_delete":
		log.Infof("[重命名] 链接「%s」->「%s」", task.Src, task.Dst)
		err = utils.CreateLink(task.Src, task.Dst)
		if err != nil {
			log.DebugErr(err)
			task.RenameState = RenameStateSeeding
			return errors.WithStack(exceptions.ErrRenameStep{Src: task.Src, Step: "链接", Message: "创建文件链接失败"})
		}
		task.RenameState = RenameStateComplete
	case "link":
		log.Infof("[重命名] 链接「%s」->「%s」", task.Src, task.Dst)
		err = utils.CreateLink(task.Src, task.Dst)
		if err != nil {
			log.DebugErr(err)
			task.RenameState = RenameStateSeeding
			return errors.WithStack(exceptions.ErrRenameStep{Src: task.Src, Step: "链接", Message: "创建文件链接失败"})
		}
		task.RenameState = RenameStateEnd
	case "move":
		log.Infof("[重命名] 移动「%s」->「%s」", task.Src, task.Dst)
		err = utils.Rename(task.Src, task.Dst)
		if err != nil {
			log.DebugErr(err)
			task.RenameState = RenameStateSeeding
			return errors.WithStack(exceptions.ErrRenameStep{Src: task.Src, Step: "移动", Message: "重命名文件失败"})
		}
		task.RenameState = RenameStateEnd
	default:
		task.RenameState = RenameStateEnd
		task.ErrCount = RenameMaxErrCount
		return errors.WithStack(exceptions.ErrRename{Src: task.Src, Message: "不支持的重命名模式 " + task.Mode})
	}
	return nil
}

func (m *Manager) stateComplete(task *RenameTask) (err error) {
	if task.ErrCount >= RenameMaxErrCount {
		log.Warnf("[重命名] 失败，跳过流程:「%s」->「%s」", task.Src, task.Dst)
		task.RenameState = RenameStateEnd
		return nil
	}
	defer func() {
		if err != nil {
			task.ErrCount++
		}
	}()
	switch task.Mode {
	case "wait_move":
		log.Infof("[重命名] 移动「%s」->「%s」", task.Src, task.Dst)
		err = utils.Rename(task.Src, task.Dst)
		if err != nil {
			log.DebugErr(err)
			task.RenameState = RenameStateComplete
			return errors.WithStack(exceptions.ErrRenameStep{Src: task.Src, Step: "移动", Message: "重命名文件失败"})
		}
		task.RenameState = RenameStateEnd
	case "link_delete":
		if !utils.IsExist(task.Dst) {
			// 确保已经链接
			err = m.stateSeeding(task)
			if err != nil {
				task.RenameState = RenameStateSeeding
				return err
			}
		}
		log.Infof("[重命名] 删除「%s」", task.Src)
		err = utils.Remove(task.Src)
		if err != nil {
			log.DebugErr(err)
			task.RenameState = RenameStateComplete
			return errors.WithStack(exceptions.ErrRenameStep{Src: task.Src, Step: "删除", Message: "删除文件失败"})
		}
		task.RenameState = RenameStateEnd
	case "link":
	case "move":
	default:
		task.RenameState = RenameStateEnd
		task.ErrCount = RenameMaxErrCount
		return errors.WithStack(exceptions.ErrRename{Src: task.Src, Message: "不支持的重命名模式 " + task.Mode})
	}
	return nil
}

func (m *Manager) AddRenameTask(opt *models.RenameOptions) (err error) {
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
			result, err = m.plugin.Rename(opt.Entity, i, xpath.Base(srcFiles[i]))
			if err != nil {
				return err
			}
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
			StateChan:      make(chan models.TorrentState, RenameStateChanCap),
			RenameCallback: opt.RenameCallback,
			RenameResult:   result,
			RenameState:    RenameStateStart,
		}
	}
	return nil
}

func (m *Manager) Update(ctx context.Context) (err error) {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	var deleteKeys []string
	for name, taskGroup := range m.taskGroups {
		if taskGroup.Complete() == AllRenameStateComplete {
			taskGroup.CompleteCallback(nil)
			deleteKeys = append(deleteKeys, name)
			continue
		}
		for _, task := range taskGroup.Tasks {
			select {
			case state := <-task.StateChan:
				task.State = state
			default:
			}
			// 初始状态
			if task.RenameState == RenameStateStart {
				if task.State != downloader.StateSeeding && task.State != downloader.StateComplete {
					continue
				}
				existSrc := utils.IsExist(task.Src)
				existDst := utils.IsExist(task.Dst)
				switch {
				case existSrc && existDst:
					// 待移动文件和目标文件都存在，覆盖
					log.Warnf("[重命名] 可能已经移动完成，覆盖:「%s」->「%s」", task.Src, task.Dst)
					task.RenameState = RenameStateSeeding
				case existSrc && !existDst:
					// 待移动文件存在，开始移动流程
					task.RenameState = RenameStateSeeding
				case !existSrc && existDst:
					// 待移动文件不存在，目标文件存在，结束移动
					log.Warnf("[重命名] 可能已经移动完成，跳过:「%s」->「%s」", task.Src, task.Dst)
					task.RenameState = RenameStateEnd
				default:
					// 待移动文件和目标文件都不存在，错误，结束移动
					return errors.WithStack(&exceptions.ErrRename{Src: task.Src, Message: "未找到文件"})
				}
			}
			// 状态一，做种
			if task.RenameState == RenameStateSeeding {
				if task.State != downloader.StateSeeding && task.State != downloader.StateComplete {
					continue
				}
				err = m.stateSeeding(task)
				if err != nil {
					return err
				}
			}
			// 状态二，完成
			if task.RenameState == RenameStateComplete {
				if task.State != downloader.StateComplete {
					continue
				}
				err = m.stateComplete(task)
				if err != nil {
					return err
				}
			}

			if task.RenameState == RenameStateEnd {
				if task.ErrCount <= RenameMaxErrCount {
					task.RenameCallback(task.RenameResult)
				}
				continue
			}

		}
	}

	for _, k := range deleteKeys {
		delete(m.taskGroups, k)
	}
	return nil
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
				var err error
				defer utils.HandleError(func(err error) {
					log.Errorf("%+v", err)
					m.sleep(ctx)
				})
				defer func() {
					if err != nil {
						log.Errorf("", err)
					}
				}()
				select {
				case <-ctx.Done():
					log.Debugf("正常退出 renamer")
					exit = true
					return
				default:
					err = m.Update(ctx)
					m.sleep(ctx)
				}
			}()
			if exit {
				return
			}
		}
	}()
}

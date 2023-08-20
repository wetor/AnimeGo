package renamer

import (
	"context"
	"path"
	"sync"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/exceptions"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/utils"
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
	Result         *models.RenameResult

	// 读写
	Enable      bool
	RenameState int
	State       models.TorrentState
	ErrCount    int
}

type RenameTaskGroup struct {
	Keys             []string
	RenameResult     *models.RenameAllResult
	CompleteCallback models.CompleteCallback // 完成重命名所有流程后回调
}

type Manager struct {
	plugin     api.RenamerPlugin
	tasks      map[string]*RenameTask
	taskGroups []*RenameTaskGroup
	sync.Mutex
}

func NewManager(plugin api.RenamerPlugin) *Manager {
	return &Manager{
		plugin:     plugin,
		tasks:      make(map[string]*RenameTask),
		taskGroups: make([]*RenameTaskGroup, 0),
	}
}

func (m *Manager) SetDownloadState(keys []string, state models.TorrentState) error {
	for _, key := range keys {
		t, ok := m.tasks[key]
		if !ok {
			return errors.WithStack(exceptions.ErrRename{Src: key, Message: "任务不存在"})
		}
		if !t.Enable {
			continue
		}
		t.StateChan <- state
	}
	return nil
}

func (m *Manager) HasRenameTask(keys []string) bool {
	for _, key := range keys {
		if _, ok := m.tasks[key]; !ok {
			return false
		}
	}
	return true
}

func (m *Manager) GetRenameTaskState(keys []string) (int, error) {
	for _, key := range keys {
		t, ok := m.tasks[key]
		if !ok {
			return AllRenameStateError, errors.WithStack(exceptions.ErrRename{Src: key, Message: "任务不存在"})
		}
		if !t.Enable {
			return AllRenameStateError, errors.WithStack(exceptions.ErrRename{Src: key, Message: "任务未启用"})
		}
	}
	state, _ := m.isComplete(keys)
	return state, nil
}

func (m *Manager) GetEpTaskState(key string) (int, error) {
	t, ok := m.tasks[key]
	if !ok {
		return RenameStateError, errors.WithStack(exceptions.ErrRename{Src: key, Message: "任务不存在"})
	}
	if !t.Enable {
		return RenameStateError, errors.WithStack(exceptions.ErrRename{Src: key, Message: "任务未启用"})
	}
	return t.RenameState, nil
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

// AddRenameTask
//
//	添加重命名任务，添加后的任务默认禁用
func (m *Manager) AddRenameTask(opt *models.RenameOptions) (renameResult *models.RenameAllResult, err error) {
	m.Lock()
	defer m.Unlock()
	renameResult = &models.RenameAllResult{
		Name:    opt.Name,
		Results: make([]*models.RenameResult, len(opt.Entity.Ep)),
	}
	srcFiles := opt.Entity.FilePathSrc()
	dstFiles := opt.Entity.FilePath()
	keys := make([]string, len(opt.Entity.Ep))
	for i := range opt.Entity.Ep {
		keys[i] = opt.Entity.EpKey(i)
		var result *models.RenameResult
		if m.plugin != nil {
			result, err = m.plugin.Rename(opt.Entity, i, path.Base(srcFiles[i]))
			if err != nil {
				return nil, err
			}
		}
		if result == nil || len(result.Filename) == 0 {
			result = &models.RenameResult{
				Index:     i,
				Scrape:    true,
				Filename:  dstFiles[i],
				AnimeDir:  opt.Entity.DirName(),
				SeasonDir: path.Dir(result.Filename),
			}
		}
		src := path.Join(opt.SrcDir, srcFiles[i])
		dst := path.Join(opt.DstDir, result.Filename)
		result.Index = i
		renameResult.Results[i] = result
		m.tasks[keys[i]] = &RenameTask{
			Enable:         false,
			Src:            src,
			Dst:            dst,
			Mode:           opt.Mode,
			StateChan:      make(chan models.TorrentState, RenameStateChanCap),
			RenameCallback: opt.RenameCallback,
			Result:         result,
			RenameState:    RenameStateStart,
		}
	}
	m.taskGroups = append(m.taskGroups, &RenameTaskGroup{
		Keys:             keys,
		RenameResult:     renameResult,
		CompleteCallback: opt.CompleteCallback,
	})
	if len(renameResult.Results) > 0 {
		renameResult.AnimeDir = renameResult.Results[0].AnimeDir
		renameResult.SeasonDir = renameResult.Results[0].SeasonDir
	}
	return renameResult, nil
}

// EnableTask
//
//	启动任务
func (m *Manager) EnableTask(keys []string) error {
	m.Lock()
	defer m.Unlock()
	for _, key := range keys {
		if task, ok := m.tasks[key]; ok {
			task.Enable = true
		} else {
			return errors.WithStack(exceptions.ErrRename{Src: key, Message: "任务不存在"})
		}
	}
	return nil
}

func (m *Manager) isComplete(keys []string) (int, bool) {
	incomplete := 0
	all := 0
	for _, key := range keys {
		if task, ok := m.tasks[key]; ok {
			if task.Enable {
				all++
				if task.RenameState != RenameStateEnd {
					incomplete++
				}
			}
		}
	}
	if all == 0 {
		return AllRenameStateError, false
	}
	if incomplete == 0 {
		return AllRenameStateComplete, true
	} else if incomplete == all {
		return AllRenameStateStart, true
	} else {
		return AllRenameStateIncomplete, true
	}
}

func (m *Manager) deleteTask(keys []string) {
	for _, key := range keys {
		delete(m.tasks, key)
	}
}

func (m *Manager) DeleteTask(keys []string) {
	m.Lock()
	defer m.Unlock()
	m.deleteTask(keys)
}

func (m *Manager) Update(ctx context.Context) (err error) {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	m.Lock()
	defer m.Unlock()

	for _, task := range m.tasks {
		if !task.Enable {
			continue
		}
		select {
		case <-ctx.Done():
			return
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
				task.RenameCallback(task.Result)
			}
			continue
		}
	}

	var deleteIndex []int
	for i, taskGroup := range m.taskGroups {
		if state, exist := m.isComplete(taskGroup.Keys); exist && state == AllRenameStateComplete {
			taskGroup.CompleteCallback(taskGroup.RenameResult)
			m.deleteTask(taskGroup.Keys)
			deleteIndex = append(deleteIndex, i)
		} else if !exist {
			m.deleteTask(taskGroup.Keys)
			deleteIndex = append(deleteIndex, i)
		}
	}
	// 清除已完成的taskGroup
	for i := len(deleteIndex) - 1; i >= 0; i-- {
		m.taskGroups = append(m.taskGroups[:deleteIndex[i]], m.taskGroups[deleteIndex[i]+1:]...)
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

package qbittorrent

import (
	"context"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	ErrorAddFailed       = "add_failed"
	ErrorListFailed      = "list_failed"
	ErrorConnectedFailed = "connected_failed"
	ErrDeleteFailed      = "delete_failed"
)

var defaultUpdateList = func(m *ClientMock) {
	for _, item := range m.Name2item {
		if item.State == QbtDownloading || item.State == QbtQueuedUP {
			item.State = QbtDownloading
			item.Progress += 0.5
			if item.Progress >= 1.005 {
				item.Progress = 0
				item.State = QbtUploading
			}
		} else if item.State == QbtUploading {
			item.Progress += 0.5
			if item.Progress >= 1.005 {
				item.State = QbtCheckingUP
			}
		}
	}
}

type ClientMockOptions struct {
	DownloadPath string
	UpdateList   func(m *ClientMock)
	Ctx          context.Context
}

type ClientMock struct {
	Name2item map[string]*models.TorrentItem
	Name2hash map[string]string
	Hash2name map[string]string
	errorFlag map[string]struct{}
	Conf      ClientMockOptions
	sync.Mutex
	Ctx context.Context
}

func (m *ClientMock) MockInit(opts ClientMockOptions) {
	m.Conf = opts

	m.Name2item = make(map[string]*models.TorrentItem)
	m.Name2hash = make(map[string]string)
	m.Hash2name = make(map[string]string)
	if m.Conf.UpdateList == nil {
		m.Conf.UpdateList = defaultUpdateList
	}
	m.errorFlag = make(map[string]struct{})
	m.Ctx = opts.Ctx
}

func (m *ClientMock) MockSetError(name string, enable bool) {
	if enable {
		m.errorFlag[name] = struct{}{}
	} else {
		delete(m.errorFlag, name)
	}
}

func (m *ClientMock) MockGetError(name string) bool {
	if _, ok := m.errorFlag[name]; ok {
		return ok
	}
	return false
}

func (m *ClientMock) MockSetUpdateList(updateList func(m *ClientMock)) {
	m.Conf.UpdateList = updateList
}

func (m *ClientMock) MockAddName(name, hash string, src []string) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorAddFailed) {
		return
	}

	m.Name2hash[name] = hash
	m.Hash2name[hash] = name

	err := utils.CreateMutiDir(path.Join(m.Conf.DownloadPath, path.Dir(src[0])))
	if err != nil {
		panic(err)
	}
	for _, s := range src {
		err = os.WriteFile(path.Join(m.Conf.DownloadPath, s), []byte{}, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func (m *ClientMock) Config() *models.Config {
	return &models.Config{
		DownloadPath: m.Conf.DownloadPath,
	}
}

func (m *ClientMock) Name() string {
	return "MockClient"
}

func (m *ClientMock) State(state string) constant.TorrentState {
	switch state {
	case QbtAllocating, QbtMetaDL, QbtStalledDL,
		QbtCheckingDL, QbtCheckingResumeData, QbtQueuedDL,
		QbtForcedUP, QbtQueuedUP:
		// 若进度为100，则下载完成
		return constant.StateWaiting
	case QbtDownloading, QbtForcedDL:
		return constant.StateDownloading
	case QbtMoving:
		return constant.StateMoving
	case QbtUploading, QbtStalledUP:
		// 已下载完成
		return constant.StateSeeding
	case QbtPausedDL:
		return constant.StatePausing
	case QbtPausedUP, QbtCheckingUP:
		// 已下载完成
		return constant.StateComplete
	case QbtError, QbtMissingFiles:
		return constant.StateError
	case QbtUnknown:
		return constant.StateUnknown
	default:
		return constant.StateUnknown
	}
}

func (m *ClientMock) Connected() bool {
	if m.MockGetError(ErrorConnectedFailed) {
		return false
	}
	return true
}

func (m *ClientMock) update() {
	m.Lock()
	defer m.Unlock()
	m.Conf.UpdateList(m)
}

func (m *ClientMock) Start() {
	go func() {
		for {
			select {
			case <-m.Ctx.Done():
				return
			default:
				m.update()
				utils.Sleep(1, m.Ctx)
			}
		}
	}()
}

func (m *ClientMock) List(opt *models.ListOptions) ([]*models.TorrentItem, error) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorListFailed) {
		return nil, errors.New(ErrorListFailed)
	}

	list := make([]*models.TorrentItem, 0, len(m.Name2item))
	for _, item := range m.Name2item {
		list = append(list, item)
	}
	return list, nil
}

func (m *ClientMock) Add(opt *models.AddOptions) error {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorAddFailed) {
		return errors.New(ErrorAddFailed)
	}

	m.Name2item[opt.Name] = &models.TorrentItem{
		ContentPath: opt.Name,
		Hash:        m.Name2hash[opt.Name],
		Name:        opt.Name,
		Progress:    0.0,
		State:       QbtQueuedUP,
	}
	return nil
}

func (m *ClientMock) Pause(opt *models.PauseOptions) error {
	return nil
}

func (m *ClientMock) Delete(opt *models.DeleteOptions) error {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrDeleteFailed) {
		return errors.New(ErrDeleteFailed)
	}

	for _, hash := range opt.Hash {
		name := m.Hash2name[hash]
		delete(m.Name2item, name)
		delete(m.Name2hash, name)
		delete(m.Hash2name, hash)
	}
	return nil
}

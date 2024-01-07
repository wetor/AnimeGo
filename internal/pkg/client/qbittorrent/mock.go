package qbittorrent

import (
	"context"
	"os"
	"path"
	"sync"

	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/internal/pkg/client"
	"github.com/wetor/AnimeGo/pkg/utils"
)

const (
	QbtQueuedUP    = "queuedUP"
	QbtDownloading = "downloading"
	QbtUploading   = "uploading"
	QbtCheckingUP  = "checkingUP"

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
}

type ClientMock struct {
	Name2item map[string]*client.TorrentItem
	Name2hash map[string]string
	Hash2name map[string]string
	errorFlag map[string]struct{}
	Conf      ClientMockOptions
	sync.Mutex
}

func (m *ClientMock) MockInit(opts ClientMockOptions) {
	m.Conf = opts

	m.Name2item = make(map[string]*client.TorrentItem)
	m.Name2hash = make(map[string]string)
	m.Hash2name = make(map[string]string)
	if m.Conf.UpdateList == nil {
		m.Conf.UpdateList = defaultUpdateList
	}
	m.errorFlag = make(map[string]struct{})
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

func (m *ClientMock) Config() *client.Config {
	return &client.Config{
		DownloadPath: m.Conf.DownloadPath,
	}
}

func (m *ClientMock) Name() string {
	return "MockClient"
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

func (m *ClientMock) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				m.update()
				utils.Sleep(1, ctx)
			}
		}
	}()
}

func (m *ClientMock) List(opt *client.ListOptions) ([]*client.TorrentItem, error) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorListFailed) {
		return nil, errors.New(ErrorListFailed)
	}

	list := make([]*client.TorrentItem, 0, len(m.Name2item))
	for _, item := range m.Name2item {
		list = append(list, item)
	}
	return list, nil
}

func (m *ClientMock) Add(opt *client.AddOptions) error {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorAddFailed) {
		return errors.New(ErrorAddFailed)
	}

	m.Name2item[opt.Name] = &client.TorrentItem{
		ContentPath: opt.Name,
		Hash:        m.Name2hash[opt.Name],
		Name:        opt.Name,
		Progress:    0.0,
		State:       QbtQueuedUP,
	}
	return nil
}

func (m *ClientMock) Delete(opt *client.DeleteOptions) error {
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

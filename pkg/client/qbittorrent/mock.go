package qbittorrent

import (
	"context"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	QbtDownloading = "downloading"
	QbtUploading   = "uploading"
	QbtCheckingUP  = "checkingUP"

	ErrorAddFailed       = "add_failed"
	ErrorListFailed      = "list_failed"
	ErrorConnectedFailed = "connected_failed"
	ErrDeleteFailed      = "delete_failed"
)

var defaultUpdateList = func(m *ClientMock) {
	for _, item := range m.name2item {
		if item.State == QbtDownloading {
			item.Progress += 0.25
			if item.Progress >= 0.5 {
				item.State = QbtUploading
			}
		} else if item.State == QbtUploading {
			item.Progress += 0.25
			if item.Progress >= 1 {
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
	name2item map[string]*client.TorrentItem
	name2hash map[string]string
	hash2name map[string]string
	errorFlag map[string]struct{}
	conf      ClientMockOptions
	sync.Mutex
}

func (m *ClientMock) MockInit(opts ClientMockOptions) {
	m.conf = opts

	m.name2item = make(map[string]*client.TorrentItem)
	m.name2hash = make(map[string]string)
	m.hash2name = make(map[string]string)
	if m.conf.UpdateList == nil {
		m.conf.UpdateList = defaultUpdateList
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
	m.conf.UpdateList = updateList
}

func (m *ClientMock) MockAddName(name, hash string, src []string) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorAddFailed) {
		return
	}

	m.name2hash[name] = hash
	m.hash2name[hash] = name

	err := utils.CreateMutiDir(xpath.Join(m.conf.DownloadPath, xpath.Dir(src[0])))
	if err != nil {
		panic(err)
	}
	for _, s := range src {
		err = os.WriteFile(xpath.Join(m.conf.DownloadPath, s), []byte{}, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
}

func (m *ClientMock) Config() *client.Config {
	return &client.Config{
		DownloadPath: m.conf.DownloadPath,
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
	m.conf.UpdateList(m)
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

	list := make([]*client.TorrentItem, 0, len(m.name2item))
	for _, item := range m.name2item {
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

	m.name2item[opt.Rename] = &client.TorrentItem{
		ContentPath: opt.Rename,
		Hash:        m.name2hash[opt.Rename],
		Name:        opt.Rename,
		Progress:    0.0,
		State:       QbtDownloading,
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
		name := m.hash2name[hash]
		delete(m.name2item, name)
		delete(m.name2hash, name)
		delete(m.hash2name, hash)
	}
	return nil
}

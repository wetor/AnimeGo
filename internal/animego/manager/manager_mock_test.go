package manager_test

import (
	"context"
	"os"
	"path"
	"sync"

	"github.com/wetor/AnimeGo/internal/models"
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

type ClientMock struct {
	name2item  map[string]*models.TorrentItem
	name2hash  map[string]string
	hash2name  map[string]string
	updateList func(m *ClientMock)
	errorFlag  map[string]struct{}
	sync.Mutex
}

func (m *ClientMock) MockInit(updateList func(m *ClientMock)) {
	m.name2item = make(map[string]*models.TorrentItem)
	m.name2hash = make(map[string]string)
	m.hash2name = make(map[string]string)
	if updateList == nil {
		updateList = defaultUpdateList
	}
	m.updateList = updateList
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
	m.updateList = updateList
}

func (m *ClientMock) MockAddName(name, hash string, src []string) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorAddFailed) {
		return
	}

	m.name2hash[name] = hash
	m.hash2name[hash] = name

	err := utils.CreateMutiDir(xpath.Join(DownloadPath, path.Dir(src[0])))
	if err != nil {
		panic(err)
	}
	for _, s := range src {
		err = os.WriteFile(xpath.Join(DownloadPath, s), []byte{}, os.ModePerm)
		if err != nil {
			panic(err)
		}
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
	m.updateList(m)
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

func (m *ClientMock) List(opt *models.ClientListOptions) []*models.TorrentItem {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorListFailed) {
		return nil
	}

	list := make([]*models.TorrentItem, 0, len(m.name2item))
	for _, item := range m.name2item {
		list = append(list, item)
	}
	return list
}

func (m *ClientMock) Add(opt *models.ClientAddOptions) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrorAddFailed) {
		return
	}

	m.name2item[opt.Rename] = &models.TorrentItem{
		ContentPath: opt.Rename,
		Hash:        m.name2hash[opt.Rename],
		Name:        opt.Rename,
		Progress:    0.0,
		State:       QbtDownloading,
	}
}

func (m *ClientMock) Delete(opt *models.ClientDeleteOptions) {
	m.Lock()
	defer m.Unlock()

	if m.MockGetError(ErrDeleteFailed) {
		return
	}

	for _, hash := range opt.Hash {
		name := m.hash2name[hash]
		delete(m.name2item, name)
		delete(m.name2hash, name)
		delete(m.hash2name, hash)
	}

}

func (m *ClientMock) GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem {

	return []*models.TorrentContentItem{}
}

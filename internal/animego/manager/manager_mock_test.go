package manager_test

import (
	"context"
	"os"
	"sync"

	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	QbtDownloading = "downloading"
	QbtUploading   = "uploading"
	QbtCheckingUP  = "checkingUP"
)

type ClientMock struct {
	name2content map[string]*models.TorrentContentItem
	name2item    map[string]*models.TorrentItem
	name2hash    map[string]string
	hash2name    map[string]string
	sync.Mutex
}

func (m *ClientMock) Init() {
	m.name2content = make(map[string]*models.TorrentContentItem)
	m.name2item = make(map[string]*models.TorrentItem)
	m.name2hash = make(map[string]string)
	m.hash2name = make(map[string]string)
}

func (m *ClientMock) AddName(name, hash string) {
	m.Lock()
	defer m.Unlock()
	m.name2hash[name] = hash
	m.hash2name[hash] = name
}

func (m *ClientMock) Connected() bool {
	return qbtConnect
}

func (m *ClientMock) update() {
	m.Lock()
	defer m.Unlock()
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
	list := make([]*models.TorrentItem, 0, len(m.name2item))
	for _, item := range m.name2item {
		list = append(list, item)
	}
	return list
}

func (m *ClientMock) Add(opt *models.ClientAddOptions) {
	m.Lock()
	defer m.Unlock()
	m.name2item[opt.Rename] = &models.TorrentItem{
		ContentPath: opt.Rename,
		Hash:        m.name2hash[opt.Rename],
		Name:        opt.Rename,
		Progress:    0.0,
		State:       QbtDownloading,
	}
	m.name2content[opt.Rename] = &models.TorrentContentItem{
		Name: xpath.Join(opt.Rename, ContentFile),
		Size: 1024,
	}

	err := utils.CreateMutiDir(xpath.Join(DownloadPath, opt.Rename))
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(xpath.Join(DownloadPath, opt.Rename, ContentFile), []byte{}, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (m *ClientMock) Delete(opt *models.ClientDeleteOptions) {
	m.Lock()
	defer m.Unlock()
	for _, hash := range opt.Hash {
		name := m.hash2name[hash]
		delete(m.name2item, name)
		delete(m.name2content, name)
		delete(m.name2hash, name)
		delete(m.hash2name, hash)
	}

}

func (m *ClientMock) GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem {
	if name, has := m.hash2name[opt.Hash]; has {
		return []*models.TorrentContentItem{m.name2content[name]}
	}
	return []*models.TorrentContentItem{}
}

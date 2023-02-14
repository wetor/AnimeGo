package downloader_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/stretchr/testify/mock"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/utils"
)

const (
	QbtDownloading = "downloading"
	QbtUploading   = "uploading"
	QbtCheckingUP  = "checkingUP"
)

var (
	name2hash    = make(map[string]string)
	itemList     = make([]*models.TorrentItem, 0)
	name2content = make(map[string]*models.TorrentContentItem, 0)
)

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) Connected() bool {
	return true
}

func (m *ClientMock) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for _, item := range itemList {
					if item.State == QbtDownloading {
						item.Progress += 0.25
						if item.Progress >= 1 {
							item.State = QbtUploading
						}
					} else if item.State == QbtUploading {
						item.Progress += 0.25
						if item.Progress >= 1.5 {
							item.State = QbtCheckingUP
						}
					}
				}
				utils.Sleep(1, ctx)
			}
		}
	}()
}

func (m *ClientMock) List(opt *models.ClientListOptions) []*models.TorrentItem {
	return itemList
}

func (m *ClientMock) Add(opt *models.ClientAddOptions) {
	itemList = append(itemList, &models.TorrentItem{
		ContentPath: opt.Rename,
		Hash:        name2hash[opt.Rename],
		Name:        opt.Rename,
		Progress:    0.0,
		State:       QbtDownloading,
	})
	name2content[opt.Rename] = &models.TorrentContentItem{
		Name: filepath.Join(opt.Rename, ContentFile),
		Size: 1024,
	}

	_ = utils.CreateMutiDir(filepath.Join(DownloadPath, opt.Rename))
	_ = os.WriteFile(filepath.Join(DownloadPath, opt.Rename, ContentFile), []byte{}, os.ModePerm)
}

func (m *ClientMock) Delete(opt *models.ClientDeleteOptions) {

}

func (m *ClientMock) GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem {
	for key, value := range name2hash {
		if value == opt.Hash {
			return []*models.TorrentContentItem{name2content[key]}
		}
	}
	return []*models.TorrentContentItem{}
}

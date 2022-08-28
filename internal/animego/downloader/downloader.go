// Package downloader
// @Description: 下载器包，用来调用外部下载器
package downloader

import (
	"GoBangumi/internal/models"
)

type Client interface {
	Version() string
	Preferences() *models.Preferences
	SetDefaultPreferences()
	List(opt *models.ClientListOptions) []*models.TorrentItem
	Rename(opt *models.ClientRenameOptions)
	Add(opt *models.ClientAddOptions)
	Delete(opt *models.ClientDeleteOptions)
	Get(opt *models.ClientGetOptions) *models.TorrentItem
	GetContent(opt *models.ClientGetOptions) []*models.TorrentContentItem
}

package client

import (
	"GoBangumi/models"
)

type Client interface {
	Version() string
	Preferences() *models.Preferences
	SetPreferences(pref *models.Preferences)
	List(opt *models.ClientListOptions) []*models.TorrentItem
	Rename(opt *models.ClientRenameOptions)
	Add(opt *models.ClientAddOptions)
	Delete(opt *models.ClientDeleteOptions)
	Get(opt *models.ClientGetOptions) []*models.TorrentItem
}

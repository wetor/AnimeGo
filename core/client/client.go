package client

import (
	"GoBangumi/model"
)

type Client interface {
	Version() string
	Preferences() *model.Preferences
	SetPreferences(pref *model.Preferences)
	List(opt *model.ClientListOptions) []*model.TorrentItem
	Rename(opt *model.ClientRenameOptions)
	Add(opt *model.ClientAddOptions)
	Delete(opt *model.ClientDeleteOptions)
	Get(opt *model.ClientGetOptions) []*model.TorrentItem
}

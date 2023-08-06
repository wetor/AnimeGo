package database

import (
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/dirdb"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Database struct {
	cacheDB api.Cacher
	dirDB   dirdb.DB
}

func NewDatabase(cacheDB api.Cacher, dirDB dirdb.DB) *Database {
	return &Database{
		cacheDB: cacheDB,
		dirDB:   dirDB,
	}
}

func (m *Database) OnDownloadStart(events []models.ClientEvent) {
	log.Infof("OnDownloadStart %v", events)
}

func (m *Database) OnDownloadPause(events []models.ClientEvent) {
	log.Infof("OnDownloadPause %v", events)
}

func (m *Database) OnDownloadStop(events []models.ClientEvent) {
	log.Infof("OnDownloadStop %v", events)
}

func (m *Database) OnDownloadSeeding(events []models.ClientEvent) {
	log.Infof("OnDownloadSeeding %v", events)
}

func (m *Database) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
}

func (m *Database) OnDownloadError(events []models.ClientEvent) {
	log.Infof("OnDownloadError %v", events)
}

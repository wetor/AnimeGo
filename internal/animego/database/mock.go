package database

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type DatabaseMock struct {
}

func (m *DatabaseMock) OnDownloadStart(events []models.ClientEvent) {
	log.Infof("OnDownloadStart %v", events)
}

func (m *DatabaseMock) OnDownloadPause(events []models.ClientEvent) {
	log.Infof("OnDownloadPause %v", events)
}

func (m *DatabaseMock) OnDownloadStop(events []models.ClientEvent) {
	log.Infof("OnDownloadStop %v", events)
}

func (m *DatabaseMock) OnDownloadSeeding(events []models.ClientEvent) {
	log.Infof("OnDownloadSeeding %v", events)
}

func (m *DatabaseMock) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
}

func (m *DatabaseMock) OnDownloadError(events []models.ClientEvent) {
	log.Infof("OnDownloadError %v", events)
}

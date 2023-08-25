package database

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type Mock struct {
}

func (m *Mock) OnDownloadStart(events []models.ClientEvent) {
	log.Infof("OnDownloadStart %v", events)
}

func (m *Mock) OnDownloadPause(events []models.ClientEvent) {
	log.Infof("OnDownloadPause %v", events)
}

func (m *Mock) OnDownloadStop(events []models.ClientEvent) {
	log.Infof("OnDownloadStop %v", events)
}

func (m *Mock) OnDownloadSeeding(events []models.ClientEvent) {
	log.Infof("OnDownloadSeeding %v", events)
}

func (m *Mock) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
}

func (m *Mock) OnDownloadError(events []models.ClientEvent) {
	log.Infof("OnDownloadError %v", events)
}

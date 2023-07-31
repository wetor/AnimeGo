package manager

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/log"
)

type ManagerMock struct {
}

func (m *ManagerMock) OnDownloadStart(events []models.ClientEvent) {
	log.Infof("OnDownloadStart %v", events)
}

func (m *ManagerMock) OnDownloadPause(events []models.ClientEvent) {
	log.Infof("OnDownloadPause %v", events)
}

func (m *ManagerMock) OnDownloadStop(events []models.ClientEvent) {
	log.Infof("OnDownloadStop %v", events)
}

func (m *ManagerMock) OnDownloadSeeding(events []models.ClientEvent) {
	log.Infof("OnDownloadSeeding %v", events)
}

func (m *ManagerMock) OnDownloadComplete(events []models.ClientEvent) {
	log.Infof("OnDownloadComplete %v", events)
}

func (m *ManagerMock) OnDownloadError(events []models.ClientEvent) {
	log.Infof("OnDownloadError %v", events)
}

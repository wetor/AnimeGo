package downloader_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/database"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/client/qbittorrent"
	"github.com/wetor/AnimeGo/pkg/log"
)

const (
	DownloadPath = "data/download"
	SavePath     = "data/save"

	HookTimeUnix = 100
)

var (
	qbt  *qbittorrent.ClientMock
	mgr  api.ClientNotifier
	dmgr *downloader.Manager
	out  *bytes.Buffer
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	_ = os.RemoveAll("data")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	qbt = &qbittorrent.ClientMock{}
	mgr = &database.DatabaseMock{}
	wg := sync.WaitGroup{}
	downloader.Init(&downloader.Options{
		RefreshSecond:          1,
		Category:               "AnimeGoTest",
		WG:                     &wg,
		AllowDuplicateDownload: false,
		SeedingTimeMinute:      0,
		Tag:                    "",
	})
	dmgr = downloader.NewManager(qbt, mgr)

	m.Run()
	_ = log.Close()
	_ = os.RemoveAll("data")
	fmt.Println("end")
}

func initTest() (*sync.WaitGroup, func()) {
	wg := sync.WaitGroup{}
	downloader.WG = &wg
	ctx, cancel := context.WithCancel(context.Background())
	qbt.MockInit(qbittorrent.ClientMockOptions{
		DownloadPath: DownloadPath,
		UpdateList: func(m *qbittorrent.ClientMock) {
			for name, item := range m.Name2item {
				if item.State == qbittorrent.QbtDownloading {
					item.Progress += 0.1
					if item.Progress >= 0.5 {
						item.State = qbittorrent.QbtUploading
					}
				} else if item.State == qbittorrent.QbtUploading {
					item.Progress += 0.1
					if item.Progress >= 1 {
						item.State = qbittorrent.QbtCheckingUP
					}
				}
				log.Debugf("%s: %v", name, int((item.Progress+0.005)*100))
			}
		},
	})
	qbt.Start(ctx)
	dmgr.Start(ctx)
	return &wg, cancel
}

func TestManager_UpdateList(t *testing.T) {
	wg, cancel := initTest()

	go func() {
		time.Sleep(13 * time.Second)
		cancel()
	}()
	hash := "326eb2aa240d4ebd8902489b031dc532"
	fullname := "test[第1季][第1集]"
	qbt.MockAddName(fullname, hash, []string{"test/S1/E1.mp4"})
	err := dmgr.Add(hash, &client.AddOptions{
		Name: fullname,
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
}

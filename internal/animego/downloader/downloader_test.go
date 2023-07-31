package downloader_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/downloader"
	"github.com/wetor/AnimeGo/internal/animego/manager"
	"github.com/wetor/AnimeGo/internal/api"
	"github.com/wetor/AnimeGo/pkg/client"
	"github.com/wetor/AnimeGo/pkg/client/qbittorrent"
	"github.com/wetor/AnimeGo/pkg/log"
	"os"
	"sync"
	"testing"
	"time"
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
	mgr = &manager.ManagerMock{}
	wg := sync.WaitGroup{}
	downloader.Init(&downloader.Options{
		RefreshSecond: 1,
		Category:      "AnimeGoTest",
		WG:            &wg,
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
	})
	qbt.Start(ctx)
	dmgr.Start(ctx)
	return &wg, cancel
}

func TestManager_UpdateList(t *testing.T) {
	wg, cancel := initTest()

	go func() {
		time.Sleep(5 * time.Second)
		cancel()
	}()
	hash := "326eb2aa240d4ebd8902489b031dc532"
	fullname := "test[第1季][第1集]"
	qbt.MockAddName(fullname, hash, []string{"test/S1/E1.mp4"})
	err := dmgr.Add(hash, &client.AddOptions{
		Rename: fullname,
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
}

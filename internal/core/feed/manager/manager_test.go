package manager

import (
	"GoBangumi/internal/core/anisource/mikan"
	mikanRss "GoBangumi/internal/core/feed/mikan"
	"GoBangumi/store"
	"GoBangumi/utils/logger"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()

	store.Init(nil)

	m.Run()
	fmt.Println("end")
}

func TestManager_UpdateFeed(t *testing.T) {
	rss := mikanRss.NewRss()
	mk := mikan.MikanAdapter{ThemoviedbKey: store.Config.KeyTmdb()}
	m := NewManager(rss, mk)

	exit := make(chan bool)
	ctx, cancel := context.WithCancel(context.Background())
	m.Start(ctx)

	go func() {
		time.Sleep(30 * time.Second)
		cancel()
		exit <- true
	}()

	//time.Sleep(120 * time.Second)

	<-exit
}

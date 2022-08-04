package manager

import (
	"GoBangumi/internal/core/anisource/mikan"
	mikanRss "GoBangumi/internal/core/feed/mikan"
	"GoBangumi/store"
	"GoBangumi/utils"
	"GoBangumi/utils/logger"
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
	mk := mikan.NewMikan()
	m := NewManager(rss, mk)

	exit := make(chan bool)
	m.Start(exit)

	go func() {
		time.Sleep(30 * time.Second)
		m.Exit()
	}()

	//time.Sleep(120 * time.Second)

	<-exit
}

func TestTimer(t *testing.T) {

	exit := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		fmt.Println("send over")
		exit <- true
	}()
	utils.Sleep(5, exit)
	fmt.Println("over")

	time.Sleep(10 * time.Second)

}

package utils

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestCreateLink(t *testing.T) {
	err := CreateLink("/Users/wetor/GoProjects/AnimeGo/data/animego.yaml", "/Users/wetor/GoProjects/AnimeGo/data/hardlink/animego.hardlink.yaml")
	if err != nil {
		log.Println(err)
	}
}

func TestConvertModel(t *testing.T) {
	ch := make(chan string)

	go func() {
		<-ch
		log.Println("ch read 11")
	}()
	time.Sleep(2 * time.Second)
	ch <- "11"
	log.Println("ch write 11")
	time.Sleep(1 * time.Second)
}

func TestRename(t *testing.T) {
	err := Rename("/Users/wetor/GoProjects/AnimeGo/download/incomplete/[ANi] PUI PUI 天竺鼠車車 駕訓班篇 - 10 [1080P][Baha][WEB-DL][AAC AVC][CHT].mp4",
		"/Users/wetor/GoProjects/AnimeGo/download/anime/PUI PUI 天竺鼠车车 DRIVING SCHOOL/S01/E010.mp4")
	if err != nil {
		log.Println(err)
	}
}

func TestRename2(t *testing.T) {
	err := os.Rename("/Users/wetor/GoProjects/AnimeGo/download/anime/PUI PUI 天竺鼠车车 DRIVING SCHOOL/S01/E010.mp4", "/Users/wetor/GoProjects/AnimeGo/download/incomplete/[ANi] PUI PUI 天竺鼠車車 駕訓班篇 - 10 [1080P][Baha][WEB-DL][AAC AVC][CHT].mp4")
	if err != nil {
		log.Println(err)
	}
}

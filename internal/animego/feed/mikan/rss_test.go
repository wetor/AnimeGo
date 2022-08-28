package mikan

import (
	"AnimeGo/internal/logger"
	"AnimeGo/internal/store"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/AnimeGo/data/config/conf.yaml",
	})
	m.Run()
	fmt.Println("end")
}

func TestRss_Parse(t *testing.T) {
	f := NewRss(store.Config.RssMikan().Url, store.Config.RssMikan().Name)
	items := f.Parse()
	for _, b := range items {
		fmt.Println(b)
	}
}

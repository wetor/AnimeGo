package mikan

import (
	"GoBangumi/internal/store"
	"GoBangumi/internal/utils/logger"
	"fmt"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	logger.Init()
	defer logger.Flush()
	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/GoBangumi/data/config/conf.yaml",
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

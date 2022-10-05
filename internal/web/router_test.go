package web

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/logger"
	"AnimeGo/internal/store"
	pkgAnisource "AnimeGo/pkg/anisource"
	"context"
	"fmt"
	"testing"
	"time"
)

func TestInitRouter(t *testing.T) {
	logger.Init()
	defer logger.Flush()
	store.Init(&store.InitOptions{
		ConfigFile: "/Users/wetor/GoProjects/AnimeGo/data/config/animego.yaml",
	})

	anisource.Init(&pkgAnisource.Options{
		Cache:   store.Cache,
		Proxy:   store.Config.Proxy(),
		Timeout: store.Config.HttpTimeoutSecond,
		Retry:   store.Config.HttpRetryNum,
	})
	var ctx, cancel = context.WithCancel(context.Background())
	store.WG.Add(1)
	Run(ctx)

	time.Sleep(3 * time.Second)
	cancel()
	store.WG.Wait()
	fmt.Println("end")
}

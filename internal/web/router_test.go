package web

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/store"
	pkgAnisource "AnimeGo/pkg/anisource"
	"AnimeGo/test"
	"context"
	"testing"
)

func TestInitRouter(t *testing.T) {
	test.TestInit()

	anisource.Init(&pkgAnisource.Options{
		Cache:   store.Cache,
		Proxy:   store.Config.Proxy(),
		Timeout: store.Config.HttpTimeoutSecond,
		Retry:   store.Config.HttpRetryNum,
	})

	Run(context.Background())

	store.WG.Wait()

}

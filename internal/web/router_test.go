package web

import (
	"AnimeGo/internal/animego/anisource"
	"AnimeGo/internal/store"
	"AnimeGo/internal/utils"
	pkgAnisource "AnimeGo/pkg/anisource"
	"AnimeGo/test"
	"context"
	"encoding/base64"
	"fmt"
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

func TestSha256(t *testing.T) {
	aa := utils.Sha256("test111222333")
	fmt.Println(aa)
	jsonStr := `
{
    "name": "filter/default.js",
    "data": "=="
}
`
	str := base64.StdEncoding.EncodeToString([]byte(jsonStr))
	fmt.Println(str)
}

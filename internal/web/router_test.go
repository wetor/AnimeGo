package web

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/anisource"
	"github.com/wetor/AnimeGo/internal/store"
	"github.com/wetor/AnimeGo/internal/utils"
	pkgAnisource "github.com/wetor/AnimeGo/pkg/anisource"
	"github.com/wetor/AnimeGo/test"
	"testing"
)

func TestInitRouter(t *testing.T) {
	test.TestInit()

	anisource.Init(&pkgAnisource.Options{
		Cache: store.Cache,
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

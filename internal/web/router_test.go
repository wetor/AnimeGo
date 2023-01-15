package web

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/utils"
)

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

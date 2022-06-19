package bangumi

import (
	"GoBangumi/models"
	"fmt"
	"testing"
)

func TestNewBgm(t *testing.T) {
	bgm := NewBgm()
	b := bgm.Parse(&models.BangumiParseOptions{
		ID: 317613,
	})
	fmt.Println(b)
}

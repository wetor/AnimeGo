package bangumi

import (
	"GoBangumi/model"
	"fmt"
	"testing"
)

func TestNewBgm(t *testing.T) {
	bgm := NewBgm()
	b := bgm.Parse(&model.BangumiParseOptions{
		ID: 317613,
	})
	fmt.Println(b)
}

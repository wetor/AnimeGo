package bangumi

import (
	"GoBangumi/models"
	"fmt"
	"testing"
)

func TestNewBgm(t *testing.T) {
	bgm := NewBgm()
	b := bgm.Parse(&models.BangumiParseOptions{
		ID:   317613,
		Ep:   6,
		Date: "2022-05-13",
	})
	fmt.Println(b, b.BangumiEp)
}

func TestBgm_Parse2(t *testing.T) {
	bgm := Bgm{}
	ep := bgm.parseBgm2(324295, 6, "2022-05-30")
	fmt.Println(ep)
}

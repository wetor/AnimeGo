package bangumi

import (
	"GoBangumi/internal/models"
	"fmt"
	"testing"
)

func TestNewBgm(t *testing.T) {
	bangumi := NewBangumi()
	b := bangumi.Parse(&models.AnimeParseOptions{
		ID:   317613,
		Ep:   6,
		Date: "2022-05-13",
	})
	fmt.Println(b, b.AnimeEp)
}

func TestBgm_Parse2(t *testing.T) {
	bangumi := Bangumi{}
	ep := bangumi.parseBnagumi2(324295, 6, "2022-05-30")
	fmt.Println(ep)
}

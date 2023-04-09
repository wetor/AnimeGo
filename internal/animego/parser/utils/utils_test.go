package utils_test

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/animego/parser/utils"
	"testing"
)

func TestParse(t *testing.T) {
	r := utils.Parse("[猎户不鸽压制] 万事屋斋藤先生转生异世界 / 斋藤先生无所不能 Benriya Saitou-san, Isekai ni Iku [01-12] [合集] [WebRip 1080p] [繁中内嵌] [H265 AAC] [2023年1月番]")
	fmt.Println(r)
}

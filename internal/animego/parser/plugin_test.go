package parser_test

import (
	"encoding/json"
	"fmt"
	parserPlugin "github.com/wetor/AnimeGo/internal/animego/parser"
	"testing"

	"github.com/wetor/AnimeGo/internal/models"
)

func TestBuiltinParser_Parse(t *testing.T) {
	p := parserPlugin.NewParserPlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_parser.py",
	}, true)

	r, _ := p.Parse("[ANi] 总之就是非常可爱 第二季（仅限港澳台地区） - 01 [1080P][Bilibili][WEB-DL][AAC AVC][CHT CHS][MP4]")
	t1 := "[ANi] 吸血鬼馬上死 第二季 - 12 [1080P][Baha][WEB-DL][AAC AVC][CHT].mp4"
	r1, _ := p.Parse(t1)
	t2 := "[orion origin] Kyuuketsuki Sugu Shinu S2 [12] [END] [1080p] [H265 AAC] [CHS].mp4"
	r2, _ := p.Parse(t2)
	d, _ := json.Marshal(r)
	fmt.Println(string(d))
	d, _ = json.Marshal(r1)
	fmt.Println(string(d))
	d, _ = json.Marshal(r2)
	fmt.Println(string(d))
}

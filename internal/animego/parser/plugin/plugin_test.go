package plugin_test

import (
	"fmt"
	parserPlugin "github.com/wetor/AnimeGo/internal/animego/parser/plugin"
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/internal/plugin"
	"github.com/wetor/AnimeGo/pkg/log"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("begin")
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})
	plugin.Init(&plugin.Options{
		Path:  "testdata",
		Debug: true,
	})
	m.Run()
	fmt.Println("end")
}

func TestParser_Parse(t *testing.T) {
	p := parserPlugin.NewParserPlugin(&models.Plugin{
		Enable: true,
		Type:   "python",
		File:   "parser.py",
	}, true)

	r := p.Parse("[ANi] 总之就是非常可爱 第二季（仅限港澳台地区） - 01 [1080P][Bilibili][WEB-DL][AAC AVC][CHT CHS][MP4]")
	fmt.Println(r)
}

func TestBuiltinParser_Parse(t *testing.T) {
	p := parserPlugin.NewParserPlugin(&models.Plugin{
		Enable: true,
		Type:   "builtin",
		File:   "builtin_parser.py",
	}, true)

	r := p.Parse("[ANi] 总之就是非常可爱 第二季（仅限港澳台地区） - 01 [1080P][Bilibili][WEB-DL][AAC AVC][CHT CHS][MP4]")
	fmt.Println(r)
}

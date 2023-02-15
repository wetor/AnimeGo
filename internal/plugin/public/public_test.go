package public_test

import (
	"fmt"
	"testing"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/internal/plugin/public"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/third_party/gpython"
)

func TestParserName(t *testing.T) {
	constant.PluginPath = "../../../assets/plugin"
	log.Init(&log.Options{
		File:  "data/log.log",
		Debug: true,
	})

	gpython.Init()
	ep := public.ParserName("【百冬练习组】【身为女主角 ～被讨厌的女主角和秘密的工作～_Heroine Tarumono!】[07][1080p AVC AAC][繁体]")
	///Users/wetor/GoProjects/AnimeGo/data/plugin/lib/Auto_Bangumi/raw_parser.py
	marshal, _ := json.Marshal(ep)

	fmt.Println(string(marshal))
}

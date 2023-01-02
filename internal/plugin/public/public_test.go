package public

import (
	"encoding/json"
	"fmt"
	"github.com/wetor/AnimeGo/test"
	"github.com/wetor/AnimeGo/third_party/gpython"
	"testing"
)

func TestParserName(t *testing.T) {
	test.TestInit()
	gpython.Init()
	ep := ParserName("【百冬练习组】【身为女主角 ～被讨厌的女主角和秘密的工作～_Heroine Tarumono!】[07][1080p AVC AAC][繁体]")
	///Users/wetor/GoProjects/AnimeGo/data/plugin/anisource/Auto_Bangumi/raw_parser.py
	marshal, _ := json.Marshal(ep)

	fmt.Println(string(marshal))
}

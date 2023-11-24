package plugin_test

import (
	"fmt"
	"github.com/wetor/AnimeGo/internal/plugin"
	"os"
	"testing"

	"github.com/wetor/AnimeGo/assets"
	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/log"
)

func TestParserName(t *testing.T) {
	plugin.Init(&plugin.Options{
		Path:  assets.TestPluginPath(),
		Debug: true,
	})
	ep, _ := plugin.ParserName("【百冬练习组】【身为女主角 ～被讨厌的女主角和秘密的工作～_Heroine Tarumono!】[07][1080p AVC AAC][繁体]")
	///Users/wetor/GoProjects/AnimeGo/data/plugin/lib/Auto_Bangumi/raw_parser.py
	marshal, _ := json.Marshal(ep)

	fmt.Println(string(marshal))
	_ = log.Close()
	_ = os.RemoveAll("data")
}

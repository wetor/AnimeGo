package gpython

import (
	_ "github.com/wetor/AnimeGo/third_party/gpython/py"
	"github.com/wetor/AnimeGo/third_party/gpython/stdlib/builtin"
	"github.com/wetor/AnimeGo/third_party/gpython/stdlib/re"
)

var isInit = false

func Init() {
	if !isInit {
		builtin.Init()
		re.Init()
		Hook()
		isInit = true
	}
}

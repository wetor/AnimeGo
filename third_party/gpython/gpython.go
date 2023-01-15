package gpython

import (
	_ "github.com/wetor/AnimeGo/third_party/gpython/py"
	"github.com/wetor/AnimeGo/third_party/gpython/stdlib/builtin"
	"github.com/wetor/AnimeGo/third_party/gpython/stdlib/re"
)

func Init() {
	builtin.Init()
	re.Init()
	Hook()
}

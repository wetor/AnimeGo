package gpython

import (
	_ "github.com/wetor/AnimeGo/third_party/gpython/py"
	_ "github.com/wetor/AnimeGo/third_party/gpython/stdlib/re"
)

var isInit = false

func Init() {
	if !isInit {
		isInit = true
	}
}

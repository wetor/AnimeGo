package python

import (
	"github.com/wetor/AnimeGo/third_party/gpython"
	"testing"
)

func TestRe1(t *testing.T) {
	gpython.Init()
	pyFile := "./script/raw_parser.py"
	err := RunWithFile(pyFile)
	if err != nil {
		panic(err)
	}

}

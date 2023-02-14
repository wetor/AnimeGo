package gpython_test

import (
	"fmt"
	"testing"

	"github.com/go-python/gpython/py"

	"github.com/wetor/AnimeGo/third_party/gpython"
)

func TestHook(t *testing.T) {
	gpython.Hook()
	obj, err := py.IntFromString("07", 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(int64(obj.(py.Int)))
}

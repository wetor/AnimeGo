package gpython

import (
	"fmt"
	"testing"

	"github.com/go-python/gpython/py"
)

func TestHook(t *testing.T) {
	Hook()
	obj, err := py.IntFromString("07", 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(int64(obj.(py.Int)))
}

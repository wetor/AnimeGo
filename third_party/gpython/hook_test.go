package gpython_test

import (
	"github.com/go-python/gpython/pytest"
	"github.com/wetor/AnimeGo/third_party/gpython"
	"testing"
)

func TestHook(t *testing.T) {
	gpython.Hook()
	pytest.RunTests(t, "testdata")
}

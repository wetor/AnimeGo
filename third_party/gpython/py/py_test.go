package py_test

import (
	"fmt"
	"testing"

	"github.com/go-python/gpython/pytest"
	_ "github.com/wetor/AnimeGo/third_party/gpython/py"
)

func TestPy(t *testing.T) {
	pytest.RunTests(t, "testdata")
}

func TestGo(t *testing.T) {
	fmt.Println(fmt.Sprintf("|%6s|", "a"))
	fmt.Println(fmt.Sprintf("|%-6s|", "a"))
	fmt.Println(fmt.Sprintf("|%-6d|", 40))
	fmt.Println(fmt.Sprintf("%06d", 40))
	fmt.Println(fmt.Sprintf("%06x", 45))
	fmt.Println(fmt.Sprintf("%06X", 45))
	fmt.Println(fmt.Sprintf("%#06x", 45))
	fmt.Println(fmt.Sprintf("%06o", 40))
	fmt.Println(fmt.Sprintf("%o", 40))
	fmt.Println(fmt.Sprintf("%#o", 40)) //不符合预期
	fmt.Println(fmt.Sprintf("%06b", 40))
	fmt.Println(fmt.Sprintf("%#06b", 40))

	fmt.Println(fmt.Sprintf("%.2f", 3.1415926))
	fmt.Println(fmt.Sprintf("%.2e", 3.1415926))
	fmt.Println(fmt.Sprintf("%.2g", 3.1415926))
}

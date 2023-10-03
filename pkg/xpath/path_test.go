package xpath_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"testing"
)

func TestRoot(t *testing.T) {
	assert.Equal(t, xpath.Root("/dir1/1.txt"), "dir1")
	assert.Equal(t, xpath.Root("../dir1/1.txt"), "..")
	assert.Equal(t, xpath.Root("/dir1/"), ".")
	assert.Equal(t, xpath.Root("dir1/"), ".")
	assert.Equal(t, xpath.Root("1.txt"), ".")
}

func ExampleRoot() {
	fmt.Println(xpath.Root("/dir1/1.txt"))
	fmt.Println(xpath.Root("/dir1/"))
	fmt.Println(xpath.Root("1.txt"))
	// Output:
	// dir1
	// .
	// .
}

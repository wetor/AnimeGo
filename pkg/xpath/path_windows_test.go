//go:build windows

package xpath_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wetor/AnimeGo/pkg/xpath"
	"testing"
)

func TestRootWindows(t *testing.T) {
	assert.Equal(t, xpath.Root("\\dir1\\1.txt"), "dir1")
	assert.Equal(t, xpath.Root("C:\\dir1\\1.txt"), "C:")
}

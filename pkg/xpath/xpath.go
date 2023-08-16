package xpath

import (
	"log"
	"path/filepath"
	"strings"
)

// Root
//
//	返回 Unix路径 或 相对路径 的最上层文件夹名
func Root(path string) string {
	p := strings.TrimPrefix(P(filepath.Clean(path)), "/")
	root := strings.SplitN(p, "/", 2)[0]
	if root == p {
		return "."
	}
	return root
}

func P(path string) string {
	return filepath.ToSlash(path)
}

func IsAbs(file string) bool {
	return filepath.IsAbs(file)
}

func Abs(path string) string {
	p, err := filepath.Abs(path)
	if err != nil {
		log.Println(err)
	}
	return P(p)
}

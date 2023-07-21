package xpath

import (
	"log"
	"path"
	"path/filepath"
)

func Join(elem ...string) string {
	return path.Join(elem...)
}

func Dir(path string) string {
	return filepath.Dir(path)
}

func Split(path string) (string, string) {
	return filepath.Split(path)
}

func Ext(path string) string {
	return filepath.Ext(path)
}

func Base(path string) string {
	return filepath.Base(path)
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

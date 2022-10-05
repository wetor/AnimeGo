package main

import (
	"AnimeGo/assets"
	"AnimeGo/internal/utils"
	"embed"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

var rootPath string
var replace bool

func main() {
	flag.StringVar(&rootPath, "path", "data", "程序资源/配置文件根目录")
	flag.BoolVar(&replace, "replace", false, "替换已存在资源/配置文件")
	flag.Parse()

	copyDir(assets.Plugin, "plugin", path.Join(rootPath, "plugin"), replace)
	copyDir(assets.Config, "config", path.Join(rootPath, "config"), replace)

}

func copyDir(fs embed.FS, src, dst string, replace bool) {
	files, err := fs.ReadDir(src)
	if err != nil {
		panic(err)
	}

	err = utils.CreateMutiDir(dst)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		srcPath := path.Join(src, file.Name())
		dstPath := path.Join(dst, file.Name())
		if file.IsDir() {
			copyDir(fs, srcPath, dstPath, replace)
			continue
		}
		fileContent, err := fs.ReadFile(srcPath)
		if err != nil {
			panic(err)
		}
		if !replace && utils.IsExist(dstPath) {
			fmt.Printf("文件[%s]已存在，是否替换[y(yes)/n(no)]: ", dstPath)
			if !scanYesNo() {
				continue
			}
		}
		if err := os.WriteFile(dstPath, fileContent, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func scanYesNo() bool {
	var s string
	_, err := fmt.Scanln(&s)
	if err != nil {
		panic(err)
	}
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

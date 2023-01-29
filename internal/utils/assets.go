package utils

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

func CopyDir(fs embed.FS, src, dst string, replace bool, skip bool) {
	files, err := fs.ReadDir(src)
	if err != nil {
		panic(err)
	}

	err = CreateMutiDir(dst)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		writeFile := true
		srcPath := path.Join(src, file.Name())
		dstPath := path.Join(dst, file.Name())
		if file.IsDir() {
			CopyDir(fs, srcPath, dstPath, replace, skip)
			continue
		}
		fileContent, err := fs.ReadFile(srcPath)
		if err != nil {
			panic(err)
		}
		if IsExist(dstPath) {
			if !replace {
				log.Printf("文件[%s]已存在，是否替换[y(yes)/n(no)]: ", dstPath)
				if !scanYesNo() {
					continue
				}
			}
			if skip {
				writeFile = false
			}
		}
		if writeFile {
			// 若已存在文件大小不一致则替换
			if int(FileSize(dstPath)) != len(fileContent) {
				log.Printf("文件 [%s] 大小改变，重新写入。", path.Base(dstPath))
				if err := os.WriteFile(dstPath, fileContent, os.ModePerm); err != nil {
					panic(err)
				}
			}
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

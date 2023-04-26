package assets

import (
	"embed"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/wetor/AnimeGo/pkg/utils"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	Dir = "plugin"

	// BuiltinPrefix
	//  内置插件前缀，不会写出
	BuiltinPrefix = "builtin"

	// BuiltinRawParser
	//  plugin/filter/Auto_Bangumi/raw_parser.py
	BuiltinRawParser = "raw_parser.py"
)

var (
	//go:embed plugin
	//go:embed plugin/filter/Auto_Bangumi/__init__.py
	Plugin embed.FS
	// BuiltinFile
	//  内置插件列表，会写出，但是内部调用的为内置文件
	BuiltinFile   = []string{BuiltinRawParser}
	BuiltinPlugin = make(map[string]*string)
)

func init() {
	loadBuiltinPlugin(Dir)
}

func loadBuiltinPlugin(src string) {
	files, err := Plugin.ReadDir(src)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		srcPath := xpath.Join(src, file.Name())
		if file.IsDir() {
			loadBuiltinPlugin(srcPath)
			continue
		}
		isBuiltin := false
		for _, name := range BuiltinFile {
			if file.Name() == name {
				isBuiltin = true
				break
			}
		}
		if isBuiltin || strings.HasPrefix(file.Name(), BuiltinPrefix) {
			data, err := Plugin.ReadFile(srcPath)
			if err != nil {
				panic(err)
			}
			dataStr := string(data)
			BuiltinPlugin[file.Name()] = &dataStr
		}
	}
}

func GetBuiltinPlugin(name string) *string {
	return BuiltinPlugin[name]
}

func WritePlugins(src, dst string, skip bool) {
	files, err := Plugin.ReadDir(src)
	if err != nil {
		panic(err)
	}

	err = utils.CreateMutiDir(dst)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		writeFile := true
		srcPath := xpath.Join(src, file.Name())
		dstPath := xpath.Join(dst, file.Name())
		if file.IsDir() {
			WritePlugins(srcPath, dstPath, skip)
			continue
		}
		if strings.HasPrefix(file.Name(), BuiltinPrefix) {
			continue
		}
		fileContent, err := Plugin.ReadFile(srcPath)
		if err != nil {
			panic(err)
		}
		if utils.IsExist(dstPath) && skip {
			writeFile = false
		}
		if writeFile {
			// Hash不一致则替换
			if utils.MD5File(dstPath) != utils.MD5(fileContent) {
				log.Printf("文件 [%s] 改变，重新写入。", xpath.Base(dstPath))
				if err := os.WriteFile(dstPath, fileContent, os.ModePerm); err != nil {
					panic(err)
				}
			}
		}
	}
}

func TestPluginPath() string {
	_, currFile, _, _ := runtime.Caller(0)
	dir := path.Dir(currFile)
	pluginDir := path.Join(dir, "plugin")
	return pluginDir
}

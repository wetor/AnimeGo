package utils

import (
	"github.com/wetor/AnimeGo/pkg/errors"
	"path"
)

// FindScript
//  @Description: 判断插件是否存在
//  @Description: 如ext='.js'，插件名可以忽略'.js'后缀；插件名也可以使用上层文件夹名，会自动加载文件夹内部的 'main.js' 或 'plugin.js'
//  @Description: 如设置为 'plugin/test'，会依次尝试加载 'plugin/test/main.js', 'plugin/test/plugin.js', 'plugin/test.js'
//  @param file string
//  @param ext string
//  @return string
//
func FindScript(file, ext string) string {
	if IsDir(file) {
		// 文件夹，在文件夹中寻找 main 和 plugin
		if IsExist(path.Join(file, "main"+ext)) {
			return path.Join(file, "main"+ext)
		} else if IsExist(path.Join(file, "plugin"+ext)) {
			return path.Join(file, "plugin"+ext)
		} else {
			errors.NewAniErrorf("插件文件夹中找不到 'main%s' 或 'plugin%s'", ext, ext).TryPanic()
		}
	} else if !IsExist(file) {
		// 文件不存在，尝试增加 ext 扩展名
		if IsExist(file + ext) {
			return file + ext
		} else {
			errors.NewAniError("插件文件不存在").TryPanic()
		}
	} else if path.Ext(file) != ext {
		// 文件存在，判断扩展名是否为 ext
		errors.NewAniError("插件文件格式错误").TryPanic()
	}
	return file
}

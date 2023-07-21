package utils

import (
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

// FindScript
//
//	@Description: 判断插件是否存在
//	@Description: 如ext='.py'，插件名可以忽略'.py'后缀；插件名也可以使用上层文件夹名，会自动加载文件夹内部的 'main.py'
//	@Description: 如设置为 'plugin/test'，会依次尝试加载 'plugin/test/main.py', 'plugin/test.py'
//	@param file string
//	@param ext string
//	@return string
func FindScript(file, ext string) string {
	if IsDir(file) {
		// 文件夹，在文件夹中寻找 main
		if IsExist(xpath.Join(file, "main"+ext)) {
			return xpath.Join(file, "main"+ext)
		} else {
			errors.NewAniErrorf("文件夹中找不到 'main%s' 插件", ext).TryPanic()
		}
	} else if !IsExist(file) {
		// 文件不存在，尝试增加 ext 扩展名
		if IsExist(file + ext) {
			return file + ext
		} else {
			errors.NewAniError("插件不存在").TryPanic()
		}
	} else if xpath.Ext(file) != ext {
		// 文件存在，判断扩展名是否为 ext
		errors.NewAniError("插件扩展名错误").TryPanic()
	}
	return file
}

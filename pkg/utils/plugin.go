package utils

import (
	"path"

	"github.com/pkg/errors"
)

// FindScript
//
//	@Description: 判断插件是否存在
//	@Description: 如ext='.py'，插件名可以忽略'.py'后缀；插件名也可以使用上层文件夹名，会自动加载文件夹内部的 'main.py'
//	@Description: 如设置为 'plugin/test'，会依次尝试加载 'plugin/test/main.py', 'plugin/test.py'
//	@param file string
//	@param ext string
//	@return string
func FindScript(file, ext string) (string, error) {
	if IsDir(file) {
		// 文件夹，在文件夹中寻找 main
		if IsExist(path.Join(file, "main"+ext)) {
			return path.Join(file, "main"+ext), nil
		} else {
			return "", errors.WithStack(errors.Errorf("文件夹中找不到 'main%s' 插件", ext))
		}
	} else if !IsExist(file) {
		// 文件不存在，尝试增加 ext 扩展名
		if IsExist(file + ext) {
			return file + ext, nil
		} else {
			return "", errors.WithStack(errors.New("插件不存在"))
		}
	} else if path.Ext(file) != ext {
		// 文件存在，判断扩展名是否为 ext
		return "", errors.WithStack(errors.New("插件扩展名错误"))
	}
	return file, nil
}

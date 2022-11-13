package javascript

import (
	"github.com/wetor/AnimeGo/internal/utils"
	"github.com/wetor/AnimeGo/pkg/errors"
	"path"
)

// FindScript
//  @Description: js插件的文件名列表，依次执行。路径相对于data_path
//  @Description: 插件名可以忽略'.js'后缀；插件名也可以使用上层文件夹名，会自动加载文件夹内部的 'main.js' 或 'plugin.js'
//  @Description: 如设置为 'plugin/test'，会依次尝试加载 'plugin/test/main.js', 'plugin/test/plugin.js', 'plugin/test.js'
//  @param file string
//  @return string
//  @return error
//
func FindScript(file string) (string, error) {
	if utils.IsDir(file) {
		// 文件夹，在文件夹中寻找 main.js 和 plugin.js
		if utils.IsExist(path.Join(file, "main.js")) {
			return path.Join(file, "main.js"), nil
		} else if utils.IsExist(path.Join(file, "plugin.js")) {
			return path.Join(file, "plugin.js"), nil
		} else {
			return "", errors.NewAniError("插件文件夹中找不到 'main.js' 或 'plugin.js'")
		}
	} else if !utils.IsExist(file) {
		// 文件不存在，尝试增加 .js 扩展名
		if utils.IsExist(file + ".js") {
			return file + ".js", nil
		} else {
			return "", errors.NewAniError("插件文件不存在")
		}
	} else if path.Ext(file) != ".js" {
		// 文件存在，判断扩展名是否为 .js
		return "", errors.NewAniError("插件文件格式错误")
	}
	return file, nil
}

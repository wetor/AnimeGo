package javascript

import _ "embed"

const (
	animeGoBaseFilename = "animego.js"
	funcMain            = "main"
)

var (
	// animeGoBaseJs 基础js文件，对象初始化时执行
	//go:embed lib/animego.js
	animeGoBaseJs string
	currRootPath  string // 当前插件根目录
	currName      string // 当前插件名
)

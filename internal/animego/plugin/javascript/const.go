package javascript

import _ "embed"

const (
	animeGoBaseFilename = "animego.js"
	funcMain            = "main"
)

var (
	// animeGoBaseJs 基础js文件，对象初始化时执行
	//go:embed animego.js
	animeGoBaseJs string
)

// Object js对象类型
type Object map[string]interface{}

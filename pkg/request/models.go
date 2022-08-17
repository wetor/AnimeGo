package request

type Param struct {
	Uri      string      // 必须
	Proxy    string      // 可选，使用代理
	SaveFile string      // 可选，保存到文件
	BindJson interface{} // 可选，绑定json到结构，结构指针
}

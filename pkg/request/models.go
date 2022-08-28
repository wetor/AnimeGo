package request

import "io"

type Param struct {
	Uri      string      // 必须
	Proxy    string      // 可选，使用代理
	SaveFile string      // 可选，保存到文件
	BindJson interface{} // 可选，绑定json到结构，结构指针
	Writer   io.Writer   // 可选，写入到writer中
	Retry    int         // 可选，重试次数，默认为1次
	Timeout  int         // 可选，超时时间
}

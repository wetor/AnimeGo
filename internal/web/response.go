// @Title		web.go
// @Description	HTTP响应体结构化
// @Author		zxysilent
// @Update		zxysilent

package web

// Response  format
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// page format
// Message
type page struct {
	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

const (
	stSucc    int = 200 //正常
	stFail    int = 300 //失败
	stErrIpt  int = 310 //输入数据有误
	stErrOpt  int = 320 //无数据返回
	stErrDeny int = 330 //没有权限
	stErrJwt  int = 340 //jwt未通过验证
	stErrSvr  int = 350 //服务端错误
	stExt     int = 400 //其他约定 //eg 更新 token
)

func newReply(code int, msg string, data ...interface{}) (int, Response) {
	if len(data) > 0 {
		return 200, Response{
			Code: code,
			Msg:  msg,
			Data: data[0],
		}
	}
	return 200, Response{
		Code: code,
		Msg:  msg,
	}
}

// Succ 返回一个成功标识的结果格式
func Succ(msg string, data ...interface{}) (int, Response) {
	return newReply(stSucc, msg, data...)
}

// Fail 返回一个失败标识的结果格式
func Fail(msg string, data ...interface{}) (int, Response) {
	return newReply(stFail, msg, data...)
}

// Page 返回一个带有分页数据的结果格式
func Page(msg string, items interface{}, count int) (int, Response) {
	return 200, Response{
		Code: stSucc,
		Msg:  msg,
		Data: page{
			Items: items,
			Count: count,
		},
	}
}

// ErrIpt 返回一个输入错误的结果格式
func ErrIpt(msg string, data ...interface{}) (int, Response) {
	return newReply(stErrIpt, msg, data...)
}

// ErrOpt 返回一个输出错误的结果格式
func ErrOpt(msg string, data ...interface{}) (int, Response) {
	return newReply(stErrOpt, msg, data...)
}

// ErrDeny 返回一个没有权限的结果格式
func ErrDeny(msg string, data ...interface{}) (int, Response) {
	return newReply(stErrDeny, msg, data...)
}

// ErrJwt 返回一个通过验证的结果格式
func ErrJwt(msg string, data ...interface{}) (int, Response) {
	return newReply(stErrJwt, msg, data...)
}

// ErrSvr 返回一个服务端错误的结果格式
func ErrSvr(msg string, data ...interface{}) (int, Response) {
	return newReply(stErrSvr, msg, data...)
}

// Ext 返回一个其他约定的结果格式
func Ext(msg string, data ...interface{}) (int, Response) {
	return newReply(stExt, msg, data...)
}

package request

import (
	"AnimeGo/third_party/goreq"
	"go.uber.org/zap"
	"io"
	"os"
	"time"
)

func request(method string, param *Param) error {
	req := goreq.Request{
		Method:  method,
		Uri:     param.Uri,
		Timeout: time.Duration(param.Timeout) * time.Second,
	}
	if len(param.Proxy) > 0 {
		req.Proxy = param.Proxy
	}
	resp, err := req.Do()
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if param.BindJson != nil {
		err = resp.Body.FromJsonTo(param.BindJson)
		if err != nil {
			return err
		}
	}
	if param.Writer != nil {
		_, err = io.Copy(param.Writer, resp.Body)
		if err != nil {
			return err
		}
	}
	if len(param.SaveFile) > 0 {
		all, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = os.WriteFile(param.SaveFile, all, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func Get(param *Param) (err error) {
	if param.Retry == 0 {
		param.Retry = 1
	}
	for i := 0; i < param.Retry; i++ {
		err = request("GET", param)
		if err != nil {
			zap.S().Debug(err)
			zap.S().Warnf("请求第%d次，失败", i+1)
		} else {
			break
		}
	}
	return err
}

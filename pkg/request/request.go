package request

import (
	"AnimeGo/third_party/goreq"
	"io"
	"os"
)

func Get(param *Param) error {
	// TODO: 增加重试机制
	req := goreq.Request{
		Method: "GET",
		Uri:    param.Uri,
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

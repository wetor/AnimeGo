package request

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/internal/constant"
	"github.com/wetor/AnimeGo/pkg/log"
)

var (
	conf = &Options{}
)

func (o *Options) Default() {
	if o.RetryWait < 3 {
		o.RetryWait = 3
	}
	if o.Timeout < 3 {
		o.Timeout = 3
	}
	if len(o.UserAgent) == 0 {
		o.UserAgent = constant.DefaultUserAgent
	}
}

func Init(opt *Options) {
	opt.Default()
	conf = opt
}

func request(uri string, method string, body interface{}, header map[string]string) *gorequest.SuperAgent {
	hostOpt := &HostOptions{}
	if conf.Host != nil {
		for host, opt := range conf.Host {
			if strings.HasPrefix(uri, host) {
				if len(opt.Redirect) > 0 {
					uri = strings.Replace(uri, host, opt.Redirect, 1)
				}
				hostOpt = opt
				break
			}
		}
	}

	method = strings.ToUpper(method)
	log.Infof("HTTP %s %s %+v", method, uri, body)
	retryWait := time.Duration(conf.RetryWait) * time.Second
	timeout := time.Duration(conf.Timeout) * time.Second
	allTimeout := timeout + (timeout+retryWait)*time.Duration(conf.Retry) // 最长等待时间

	var m *gorequest.SuperAgent
	switch method {
	case "GET":
		m = gorequest.New().Get(uri)
	case "POST":
		m = gorequest.New().Post(uri)
	}
	agent := m.Send(body).
		Timeout(allTimeout).
		Proxy(conf.Proxy).
		// SetDebug(conf.Debug).
		Retry(conf.Retry, retryWait,
			http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout)

	if hostOpt.Params != nil {
		agent.Query(hostOpt.Params)
	}
	if header != nil || hostOpt.Header != nil {
		for k, v := range header {
			agent.Set(strings.ToLower(k), v)
		}
		for k, v := range hostOpt.Header {
			agent.Set(strings.ToLower(k), v)
		}
	}
	if hostOpt.Cookie != nil {
		for k, v := range hostOpt.Cookie {
			v = strings.TrimPrefix(v, k+"=")
			agent.AddCookie(&http.Cookie{
				Name:  k,
				Value: v,
			})
		}
	}

	agent.Set("user-agent", conf.UserAgent)
	return agent
}

func handleError(resp gorequest.Response, errs []error) (err error) {
	if len(errs) != 0 {
		return errs[0]
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("HTTP 请求失败, %s", resp.Status)
	}
	return nil
}

// GetString
//
//	uri string 请求地址
//	args[0] map[string]string 请求头header
func GetString(uri string, args ...interface{}) (string, error) {
	var header map[string]string = nil
	if len(args) > 0 {
		header = args[0].(map[string]string)
	}
	resp, str, errs := request(uri, "GET", nil, header).End()
	err := handleError(resp, errs)
	if err != nil {
		return "", err
	}
	return str, nil
}

func Get(uri string, body interface{}, args ...interface{}) error {
	var header map[string]string = nil
	if len(args) > 0 {
		header = args[0].(map[string]string)
	}
	resp, _, errs := request(uri, "GET", nil, header).EndStruct(body)
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	return nil
}

func Post(uri string, req interface{}, body interface{}, args ...interface{}) error {
	var header map[string]string = nil
	if len(args) > 0 {
		header = args[0].(map[string]string)
	}
	resp, _, errs := request(uri, "POST", req, header).EndStruct(body)
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	return nil
}

func GetFile(uri string, file string, args ...interface{}) error {
	var header map[string]string = nil
	if len(args) > 0 {
		header = args[0].(map[string]string)
	}
	resp, bodyBytes, errs := request(uri, "GET", nil, header).EndBytes()
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, bodyBytes, constant.WriteFilePerm)
	if err != nil {
		return err
	}
	return nil
}

func GetWriter(uri string, w io.Writer, args ...interface{}) error {
	var header map[string]string = nil
	if len(args) > 0 {
		header = args[0].(map[string]string)
	}
	resp, bodyBytes, errs := request(uri, "GET", nil, header).EndBytes()
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	_, err = w.Write(bodyBytes)
	if err != nil {
		return err
	}
	return nil
}

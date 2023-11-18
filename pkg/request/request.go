package request

import (
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"

	"github.com/wetor/AnimeGo/pkg/log"
)

var (
	Retry     int
	RetryWait int
	Timeout   int
	Proxy     string
	UserAgent string
	Debug     bool
	ReInitWG  sync.WaitGroup
)

type Options struct {
	UserAgent string
	Proxy     string // 使用代理
	Retry     int    // 额外重试次数，默认为不做重试
	RetryWait int    // 重试等待时间，最小3秒
	Timeout   int    // 超时时间，最小3秒
	Debug     bool
}

func (o *Options) Default() {
	if o.RetryWait < 3 {
		o.RetryWait = 3
	}
	if o.Timeout < 3 {
		o.Timeout = 3
	}
	o.UserAgent = "0.1.0/AnimeGo (https://github.com/wetor/AnimeGo)"
}

func Init(opt *Options) {
	opt.Default()
	Retry = opt.Retry
	RetryWait = opt.RetryWait
	Timeout = opt.Timeout
	Proxy = opt.Proxy
	UserAgent = opt.UserAgent
}

func ReInit(opt *Options) {
	ReInitWG.Wait()
	Retry = opt.Retry
	RetryWait = opt.RetryWait
	Timeout = opt.Timeout
	Proxy = opt.Proxy
	UserAgent = opt.UserAgent
}

func request(uri string, method string, body interface{}, header map[string]string) *gorequest.SuperAgent {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	method = strings.ToUpper(method)
	log.Infof("HTTP %s %s %+v", method, uri, body)
	retryWait := time.Duration(RetryWait) * time.Second
	timeout := time.Duration(Timeout) * time.Second
	allTimeout := timeout + (timeout+retryWait)*time.Duration(Retry) // 最长等待时间

	var m *gorequest.SuperAgent
	switch method {
	case "GET":
		m = gorequest.New().Get(uri)
	case "POST":
		m = gorequest.New().Post(uri)
	}
	agent := m.Send(body).
		Timeout(allTimeout).
		Proxy(Proxy).
		SetDebug(Debug).
		Retry(Retry, retryWait,
			http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout)

	if header != nil {
		for k, v := range header {
			agent.Set(strings.ToLower(k), v)
		}
	}
	agent.Set("user-agent", UserAgent)
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
	ReInitWG.Add(1)
	defer ReInitWG.Done()
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
	ReInitWG.Add(1)
	defer ReInitWG.Done()
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
	ReInitWG.Add(1)
	defer ReInitWG.Done()
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
	ReInitWG.Add(1)
	defer ReInitWG.Done()
	var header map[string]string = nil
	if len(args) > 0 {
		header = args[0].(map[string]string)
	}
	resp, bodyBytes, errs := request(uri, "GET", nil, header).EndBytes()
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, bodyBytes, 0666)
	if err != nil {
		return err
	}
	return nil
}

func GetWriter(uri string, w io.Writer, args ...interface{}) error {
	ReInitWG.Add(1)
	defer ReInitWG.Done()
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

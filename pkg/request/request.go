package request

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/pkg/log"
)

var (
	Retry     int
	RetryWait int
	Timeout   int
	Proxy     string
	UserAgent string
	Debug     bool
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

func get(uri string) *gorequest.SuperAgent {
	log.Infof("HTTP GET %s", uri)
	retryWait := time.Duration(RetryWait) * time.Second
	timeout := time.Duration(Timeout) * time.Second
	allTimeout := timeout + (timeout+retryWait)*time.Duration(Retry) // 最长等待时间
	return gorequest.New().
		Timeout(allTimeout).
		Proxy(Proxy).
		SetDebug(Debug).
		Get(uri).
		Set("User-Agent", UserAgent).
		Retry(Retry, retryWait,
			http.StatusBadRequest,
			http.StatusNotFound,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable,
			http.StatusGatewayTimeout)
}

func handleError(resp gorequest.Response, errs []error) (err error) {
	if len(errs) != 0 {
		log.Debugf("", errors.NewAniErrorD(errs))
		log.Warnf("HTTP 请求失败")
		return errs[0]
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.NewAniErrorSkipf(3, "HTTP 请求失败，%s, 重试 %s 次", resp.Status, resp.Header.Get("Retry-Count"))
		log.Debugf("", err)
		log.Warnf("HTTP 请求失败, %s", resp.Status)
		return err
	}

	if retryCount := resp.Header.Get("Retry-Count"); retryCount != "0" {
		log.Infof("HTTP 请求完成，重试 %s 次", retryCount)
	} else {
		log.Infof("HTTP 请求完成")
	}
	return nil
}

func GetString(uri string) (string, error) {
	resp, str, errs := get(uri).End()
	err := handleError(resp, errs)
	if err != nil {
		return "", err
	}
	return str, nil
}

func Get(uri string, body interface{}) error {
	resp, _, errs := get(uri).EndStruct(body)
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	return nil
}

func GetFile(uri string, file string) error {
	resp, bodyBytes, errs := get(uri).EndBytes()
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	err = os.WriteFile(file, bodyBytes, 0666)
	if err != nil {
		return errors.NewAniErrorD(err)
	}
	return nil
}

func GetWriter(uri string, w io.Writer) error {
	resp, bodyBytes, errs := get(uri).EndBytes()
	err := handleError(resp, errs)
	if err != nil {
		return err
	}
	_, err = w.Write(bodyBytes)
	if err != nil {
		return errors.NewAniErrorD(err)
	}
	return nil
}

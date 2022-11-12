package request

import (
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/wetor/AnimeGo/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	Retry     int
	RetryWait int
	Timeout   int
	Proxy     string
	UserAgent string
	Debug     bool
)

func Init(opt *InitOptions) {
	opt.Default()
	Retry = opt.Retry
	RetryWait = opt.RetryWait
	Timeout = opt.Timeout
	Proxy = opt.Proxy
	UserAgent = fmt.Sprintf("%s/AnimeGo (%s)", os.Getenv("ANIMEGO_VERSION"), os.Getenv("ANIMEGO_GITHUB"))
}

func get(uri string) *gorequest.SuperAgent {
	zap.S().Infof("HTTP GET %s", uri)
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
		zap.S().Debug(errors.NewAniErrorD(errs))
		zap.S().Warn("HTTP 请求失败")
		return errs[0]
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.NewAniErrorSkipf(3, "HTTP 请求失败，%s, 重试 %s 次", nil, resp.Status, resp.Header.Get("Retry-Count"))
		zap.S().Debug(err)
		zap.S().Warnf("HTTP 请求失败, %s", resp.Status)
		return err
	}

	if retryCount := resp.Header.Get("Retry-Count"); retryCount != "0" {
		zap.S().Infof("HTTP 请求完成，重试 %s 次", retryCount)
	} else {
		zap.S().Infof("HTTP 请求完成")
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

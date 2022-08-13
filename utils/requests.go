package utils

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func GetHttpClient(proxyUrl string) (*http.Client, error) {
	var client *http.Client
	if len(proxyUrl) > 0 {
		proxy, err := url.Parse(proxyUrl)
		if err != nil {
			return nil, err
		}
		netTransport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			ResponseHeaderTimeout: time.Second * time.Duration(5),
		}
		client = &http.Client{
			Timeout:   time.Second * 10,
			Transport: netTransport,
		}
	} else {
		client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	return client, nil
}

func HttpGet(url_, savePath, proxyUrl string) error {
	client, err := GetHttpClient(proxyUrl)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, data, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func ApiGet(url_ string, obj interface{}, proxyUrl string) (int, error) {
	client, err := GetHttpClient(proxyUrl)
	if err != nil {
		return 0, err
	}
	request, err := http.NewRequest("GET", url_, nil)
	if err != nil {
		return 0, err
	}
	request.Header.Set("User-Agent", "AnimeGo/1.0 (Golang 1.18)")
	request.Header.Set("Accept", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	if obj == nil {
		return response.StatusCode, nil
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return 0, err
	}
	return response.StatusCode, nil
}

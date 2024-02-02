package request_test

import (
	"log"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wetor/AnimeGo/internal/pkg/request"
)

func TestHost(t *testing.T) {
	reqCount := 0
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		if reqCount >= 1 && reqCount < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			reqCount++
			return
		}

		assert.Equal(t, "TestHeaderValue", r.Header.Get("TestHeaderKey"))
		assert.Equal(t, "TestHeaderValue2", r.Header.Get("TestHeaderKey2"))
		assert.Equal(t, "TestParamsValue", r.FormValue(strings.ToLower("TestParamsKey")))
		assert.Equal(t, "255", r.FormValue(strings.ToLower("data")))
		c, err := r.Cookie("TestCookieKey")
		assert.Nil(t, err)
		assert.Equal(t, "TestCookieValue", c.Value)
		_, _ = w.Write([]byte("world"))
		// 第2, 3次请求固定失败
		reqCount++
	})

	log.Println("Starting server...")
	l, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		log.Fatal(http.Serve(l, nil))
	}()

	request.Init(&request.Options{
		Host: map[string]*request.HostOptions{
			"http://192.168.1.1:8080": {
				Redirect: "http://localhost:8080",
				Header: map[string]string{
					"TestHeaderKey": "TestHeaderValue",
				},
				Params: map[string]string{
					"TestParamsKey": "TestParamsValue",
				},
				Cookie: map[string]string{
					"TestCookieKey": "TestCookieKey=TestCookieValue",
				},
			},
		},
		Retry:     3,
		RetryWait: 2,
		Debug:     true,
	})

	res, err := request.GetString("http://192.168.1.1:8080/hello?data=255", map[string]string{
		"TestHeaderKey2": "TestHeaderValue2",
	})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "world", res)

	res, err = request.GetString("http://192.168.1.1:8080/hello?data=255", map[string]string{
		"TestHeaderKey2": "TestHeaderValue2",
	})
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "world", res)
}

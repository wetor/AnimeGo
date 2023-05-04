package test

import (
	"encoding/json"
	"io"

	"github.com/brahma-adshonor/gohook"

	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
	"github.com/wetor/AnimeGo/pkg/xpath"
)

const (
	Get       = "request.Get"
	GetWriter = "request.GetWriter"
	GetString = "request.GetString"
)

var (
	testdataAll string // 目录名
	filenameAll func(string) string

	testdata map[string]string              // 目录名
	filename map[string]func(string) string // 文件名函数
)

func init() {
	testdata = make(map[string]string)
	filename = make(map[string]func(string) string)
}

func HookAll(testdataDir string, filenameFunc func(string) string) {
	HookGetWriter(testdataDir, filenameFunc)
	HookGet(testdataDir, filenameFunc)
	HookGetString(testdataDir, filenameFunc)
}

func UnHook() {
	_ = gohook.UnHook(request.GetWriter)
	_ = gohook.UnHook(request.Get)
	_ = gohook.UnHook(request.GetString)
}

func HookGetWriter(testdataDir string, filenameFunc func(string) string) {
	testdata[GetWriter] = testdataDir
	if filenameFunc == nil {
		filenameFunc = xpath.Base
	}
	filename[GetWriter] = filenameFunc
	err := gohook.Hook(request.GetWriter, getWriter, nil)
	if err != nil {
		panic(err)
	}
}

func HookGet(testdataDir string, filenameFunc func(string) string) {
	testdata[Get] = testdataDir
	if filenameFunc == nil {
		filenameFunc = xpath.Base
	}
	filename[Get] = filenameFunc
	err := gohook.Hook(request.Get, get, nil)
	if err != nil {
		panic(err)
	}
}

func HookGetString(testdataDir string, filenameFunc func(string) string) {
	testdata[GetString] = testdataDir
	if filenameFunc == nil {
		filenameFunc = xpath.Base
	}
	filename[GetString] = filenameFunc
	err := gohook.Hook(request.GetString, getString, nil)
	if err != nil {
		panic(err)
	}
}

func getWriter(uri string, w io.Writer) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := filename[GetWriter](uri)
	jsonData := GetData(testdata[GetWriter], id)
	_, err := w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func get(uri string, body interface{}) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := filename[Get](uri)
	jsonData := GetData(testdata[Get], id)
	err := json.Unmarshal(jsonData, body)
	if err != nil {
		return err
	}
	return nil
}

func getString(uri string, args ...interface{}) (string, error) {
	log.Infof("Mock HTTP GET %s, header %s", uri, args)
	id := filename[GetString](uri)
	jsonData := GetData(testdata[GetString], id)
	return string(jsonData), nil
}

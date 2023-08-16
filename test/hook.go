package test

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/wetor/AnimeGo/pkg/log"
	"github.com/wetor/AnimeGo/pkg/request"
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

	patches = gomonkey.NewPatches()
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
	patches.Reset()
}

func Hook(target interface{}, replace interface{}) {
	patches.ApplyFunc(target, replace)
}

func HookSingle(target interface{}, replace interface{}) *gomonkey.Patches {
	return gomonkey.ApplyFunc(target, replace)
}

func HookGetWriter(testdataDir string, filenameFunc func(string) string) {
	testdata[GetWriter] = testdataDir
	if filenameFunc == nil {
		filenameFunc = filepath.Base
	}
	filename[GetWriter] = filenameFunc
	patches.ApplyFunc(request.GetWriter, getWriter)
}

func HookGet(testdataDir string, filenameFunc func(string) string) {
	testdata[Get] = testdataDir
	if filenameFunc == nil {
		filenameFunc = filepath.Base
	}
	filename[Get] = filenameFunc
	patches.ApplyFunc(request.Get, get)
}

func HookGetString(testdataDir string, filenameFunc func(string) string) {
	testdata[GetString] = testdataDir
	if filenameFunc == nil {
		filenameFunc = filepath.Base
	}
	filename[GetString] = filenameFunc
	patches.ApplyFunc(request.GetString, getString)
}

func getWriter(uri string, w io.Writer) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := filename[GetWriter](uri)
	jsonData, err := GetData(testdata[GetWriter], id)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func get(uri string, body interface{}) error {
	log.Infof("Mock HTTP GET %s", uri)
	id := filename[Get](uri)
	jsonData, err := GetData(testdata[Get], id)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, body)
	if err != nil {
		return err
	}
	return nil
}

func getString(uri string, args ...interface{}) (string, error) {
	log.Infof("Mock HTTP GET %s, header %s", uri, args)
	id := filename[GetString](uri)
	jsonData, err := GetData(testdata[GetString], id)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

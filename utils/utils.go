package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"reflect"
)

func ConvertModel(src, dst interface{}) {
	vsrc := reflect.ValueOf(src).Elem()
	vscrType := vsrc.Type()
	vdst := reflect.ValueOf(dst).Elem()
	for i := 0; i < vscrType.NumField(); i++ {
		v := vdst.FieldByName(vscrType.Field(i).Name)
		if v.CanSet() {
			v.Set(vsrc.Field(i))
		}
	}
}
func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func HttpGet(url, savePath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, all, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

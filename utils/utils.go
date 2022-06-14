package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"time"
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
func CompareTime(t1, t2 time.Time, day int) bool {
	ut1 := t1.Unix()
	ut2 := t2.Unix()
	if ut1 <= ut2 {
		return ut2-ut1 <= int64(day*24*60*60)
	} else {
		return ut1-ut2 <= int64(day*24*60*60)
	}
}
func CompareTimeStr(t1, t2 string, day int) bool {
	time1, _ := time.Parse("2006-01-02", t1)
	time2, _ := time.Parse("2006-01-02", t2)
	ut1 := time1.Unix()
	ut2 := time2.Unix()
	if ut1 <= ut2 {
		return ut2-ut1 <= int64(day*24*60*60)
	} else {
		return ut1-ut2 <= int64(day*24*60*60)
	}
}
func StrTimeSub(t1, t2 string) int {
	time1, _ := time.Parse("2006-01-02", t1)
	time2, _ := time.Parse("2006-01-02", t2)
	ut1 := time1.Unix()
	ut2 := time2.Unix()
	if ut1 <= ut2 {
		return int(ut2-ut1) / (24 * 60 * 60)
	} else {
		return int(ut1-ut2) / (24 * 60 * 60)
	}
}
func GetTimeRangeDay(t time.Time, day int) (lower, upper string) {
	lower = t.Add(-time.Duration(day) * time.Hour * 24).Format("2006-01-02")
	upper = t.Add(time.Duration(day) * time.Hour * 24).Format("2006-01-02")
	return lower, upper
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
func ApiGet(url string, obj interface{}) (int, error) {

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}
	request.Header.Set("User-Agent", "GoBangumi/1.0 (Golang 1.18)")
	request.Header.Set("Accept", "application/json")
	response, _ := client.Do(request)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(data, obj)
	if err != nil {
		return 0, err
	}
	status := response.StatusCode
	return status, nil
}

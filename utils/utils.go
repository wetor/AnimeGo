package utils

import (
	"crypto/md5"
	"encoding/hex"
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
func StrTimeSubAbs(t1, t2 string) int {
	s := StrTimeSub(t1, t2)
	if s < 0 {
		return -int(s / (24 * 60 * 60))
	} else {
		return int(s / (24 * 60 * 60))
	}
}
func StrTimeSub(t1, t2 string) int64 {
	time1, _ := time.Parse("2006-01-02", t1)
	time2, _ := time.Parse("2006-01-02", t2)
	ut1 := time1.Unix()
	ut2 := time2.Unix()
	return ut1 - ut2
}
func GetTimeRangeDay(t time.Time, day int) (lower, upper string) {
	lower = t.Add(-time.Duration(day) * time.Hour * 24).Format("2006-01-02")
	upper = t.Add(time.Duration(day) * time.Hour * 24).Format("2006-01-02")
	return lower, upper
}

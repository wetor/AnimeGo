package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type FormatMap map[string]interface{}

func Format(format string, p FormatMap) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

var filenameMap = FormatMap{
	`/`: "",
	`\`: "",
	`[`: "(",
	`]`: ")",
	`:`: "-",
	`;`: "-",
	`=`: "-",
	`,`: "-",
}

func Filename(filename string) string {
	return Format(filename, filenameMap)
}

// Sleep
//  @Description: 信号计时器，每秒检测一次信号，避免长时间等待无法接收信号
//  @Description: 收到exit信号后，会返回true；倒计时结束，会返回false
//  @param second int
//  @param exit chan bool
//  @return bool
//
func Sleep(second int, ctx context.Context) bool {
	for second > 0 {
		select {
		case <-ctx.Done():
			return true
		default:
			second--
			time.Sleep(time.Second)
		}
	}
	return false
}

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

//CreateMutiDir 调用os.MkdirAll递归创建文件夹
func CreateMutiDir(filePath string) error {
	if !isExist(filePath) {
		return os.MkdirAll(filePath, os.ModePerm)
	}
	return nil
}

//isExist 判断所给路径文件/文件夹是否存在(返回true是存在)
func isExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
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
	if len(t1)*len(t2) == 0 {
		return 0
	}
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

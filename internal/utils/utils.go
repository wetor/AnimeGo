package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func Format(format string, p map[string]any) string {
	args, i := make([]string, len(p)*2), 0
	for k, v := range p {
		args[i] = "{" + k + "}"
		args[i+1] = fmt.Sprint(v)
		i += 2
	}
	return strings.NewReplacer(args...).Replace(format)
}

func Tag(tagSrc string, airDate string, ep int) string {
	date, _ := time.Parse("2006-01-02", airDate)
	mouth := (int(date.Month()) + 2) / 3
	tag := Format(tagSrc, map[string]any{
		"year":          date.Year(),
		"quarter":       (mouth-1)*3 + 1,
		"quarter_index": mouth,
		"quarter_name":  []string{"冬", "春", "夏", "秋"}[mouth-1],
		"ep":            ep,
		"week":          (int(date.Weekday())+6)%7 + 1,
		"week_name":     []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}[date.Weekday()],
	})
	return tag
}

// Sleep
//
//	@Description: 信号计时器，每秒检测一次信号，避免长时间等待无法接收信号
//	@Description: 收到exit信号后，会返回true；倒计时结束，会返回false
//	@param second int
//	@param exit chan bool
//	@return bool
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

func Map2Model(src map[string]any, dst any) {
	vdst := reflect.ValueOf(dst).Elem()
	for key, val := range src {
		v := vdst.FieldByName(key)
		if v.CanSet() {
			v.Set(reflect.ValueOf(val))
		}
	}
}

func Model2Map(src any, dst map[string]any) {
	vsrcp := reflect.ValueOf(src)
	if vsrcp.IsNil() {
		return
	}
	vsrc := vsrcp.Elem()
	vscrType := vsrc.Type()
	for i := 0; i < vscrType.NumField(); i++ {
		dst[vscrType.Field(i).Name] = vsrc.Field(i).Interface()
	}
}

func Map2ModelByJson(src map[string]any, dst any) {
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, dst)
	if err != nil {
		panic(err)
	}
}

func Model2MapByJson(src any, dst map[string]any) {
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &dst)
	if err != nil {
		panic(err)
	}
}

func ConvertModel(src, dst any) {
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

// CreateMutiDir 调用os.MkdirAll递归创建文件夹
func CreateMutiDir(filePath string) error {
	if !IsExist(filePath) {
		return os.MkdirAll(filePath, os.ModePerm)
	}
	return nil
}

// IsExist 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// FileSize 获取文件大小，文件不存返回-1
func FileSize(path string) int64 {
	s, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return s.Size()
		}
		return -1
	}
	return s.Size()
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func Md5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

func UTCToTimeStr(t1 string) string {
	time1, _ := time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", t1)
	return time1.Format("2006-01-02")
}

func CreateLink(src, dst string) error {
	dir := filepath.Dir(dst)
	if err := CreateMutiDir(dir); err != nil {
		return err
	}
	if IsExist(dst) {
		if err := os.Remove(dst); err != nil {
			return err
		}
	}
	if err := os.Link(src, dst); err != nil {
		return err
	}
	return nil
}

func Rename(src, dst string) error {
	dir := filepath.Dir(dst)
	if err := CreateMutiDir(dir); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

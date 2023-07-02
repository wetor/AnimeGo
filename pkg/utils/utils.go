package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/wetor/AnimeGo/pkg/json"
	"github.com/wetor/AnimeGo/pkg/xpath"
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

func MapToStruct(src map[string]any, dst any) {
	data, err := json.Marshal(src)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, dst)
	if err != nil {
		panic(err)
	}
}

func StructToMap(src any) (dst map[string]any) {
	vsrcp := reflect.ValueOf(src)
	var value any
	var vsrc reflect.Value

	if vsrcp.Type().Kind() == reflect.Pointer {
		if vsrcp.IsNil() {
			return
		}
		vsrc = vsrcp.Elem()
	} else {
		vsrc = vsrcp
	}

	dst = make(map[string]any)
	vscrType := vsrc.Type()
	for i := 0; i < vscrType.NumField(); i++ {
		field := vscrType.Field(i)
		value = vsrc.Field(i).Interface()
		switch field.Type.Kind() {
		case reflect.Struct:
			fallthrough
		case reflect.Pointer:
			value = StructToMap(value)
		}
		keyName := field.Tag.Get("json")
		if len(keyName) == 0 {
			keyName = field.Name
		}
		dst[keyName] = value
	}
	return dst
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

func MD5Str(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func MD5(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func MD5File(file string) string {
	f, _ := os.ReadFile(file)
	h := md5.New()
	h.Write(f)
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
	dir := xpath.Dir(dst)
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
	dir := xpath.Dir(dst)
	if err := CreateMutiDir(dir); err != nil {
		return err
	}

	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}

func Unix() int64 {
	return time.Now().Unix()
}

func String2Bool(str string) bool {
	switch strings.ToLower(str) {
	case "true", "1":
		return true
	default:
		return false
	}
}

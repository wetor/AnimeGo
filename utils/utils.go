package utils

import (
	"GoBangumi/models"
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/gob"
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

// ToBytes
//  @Description:
//  @param val
//  @param extra 若为-1则仅用作key，无法转换为value
//  @return []byte
//
func ToBytes(val interface{}, extra int64) []byte {
	buf := bytes.NewBuffer(nil)
	if extra >= 0 {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(extra))
		buf.Write(b)
	}
	switch value := val.(type) {

	case bool:
		if !value {
			buf.WriteByte(0x00)
		} else {
			buf.WriteByte(0x01)
		}
	case int:
		b := make([]byte, 4)
		binary.LittleEndian.PutUint32(b, uint32(value))
		buf.WriteByte(0x04)
		buf.Write(b)
	case int64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(value))
		buf.WriteByte(0x08)
		buf.Write(b)
	case string:
		buf.WriteByte(0x10)
		buf.WriteString(value)
	case *models.Bangumi:
		buf.WriteByte(0x20)
		buf.Write(GobToBytes(value))
	case *models.BangumiSeason:
		buf.WriteByte(0x21)
		buf.Write(GobToBytes(value))
	case *models.BangumiEp:
		buf.WriteByte(0x22)
		buf.Write(GobToBytes(value))
	case *models.BangumiExtra:
		buf.WriteByte(0x23)
		buf.Write(GobToBytes(value))
	default:
		buf.WriteByte(0xFF)
		buf.Write(GobToBytes(value))
	}
	return buf.Bytes()
}

func ToValue(data []byte) (val interface{}, extra int64) {
	_ = data[8]
	extra = int64(binary.LittleEndian.Uint64(data[0:8]))
	switch data[8] {
	case 0x00:
		val = false
	case 0x01:
		val = true
	case 0x04:
		val = int(binary.LittleEndian.Uint32(data[9:]))
	case 0x08:
		val = int64(binary.LittleEndian.Uint64(data[9:]))
	case 0x10:
		val = string(data[9:])
	case 0x20:
		val = &models.Bangumi{}
		GobToValue(data[9:], val)
	case 0x21:
		val = &models.BangumiSeason{}
		GobToValue(data[9:], val)
	case 0x22:
		val = &models.BangumiEp{}
		GobToValue(data[9:], val)
	case 0x23:
		val = &models.BangumiExtra{}
		GobToValue(data[9:], val)
	case 0xFF:
		GobToValue(data[9:], val)
	}
	return val, extra
}

func GobToBytes(val interface{}) []byte {
	buf2 := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf2)
	enc.Encode(val)
	return buf2.Bytes()
}
func GobToValue(data []byte, val interface{}) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	dec.Decode(val)
}

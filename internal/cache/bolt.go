package cache

import (
	"GoBangumi/internal/models"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"go.uber.org/zap"
	"path"
	"time"
)

type Bolt struct {
	db *bolt.DB
}

func NewBolt() Cache {
	return &Bolt{}
}

func (c *Bolt) Open(dir string) {
	db, err := bolt.Open(path.Join(dir, "bolt.db"), 0600, nil)
	if err != nil {
		zap.S().Warn(err)
		return
	}

	c.db = db

	err = c.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range buckets {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		zap.S().Warn(err)
		return
	}

}

func (c *Bolt) Close() {
	err := c.db.Close()
	if err != nil {
		zap.S().Warn(err)
		return
	}
}

func (c *Bolt) Put(bucket string, key, val interface{}, ttl int64) {
	if val == nil {
		return
	}
	var expire int64
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if ttl > 0 {
			expire = time.Now().Unix() + ttl
		} else {
			expire = 0
		}
		err := b.Put(c.toBytes(key, -1), c.toBytes(val, expire))
		return err
	})
	if err != nil {
		zap.S().Warn(err)
		return
	}
}

func (c *Bolt) Get(bucket string, key interface{}) interface{} {
	var val interface{}
	var ttl int64
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get(c.toBytes(key, -1))
		if len(v) == 0 {
			// 不存在
			return nil
		}
		val, ttl = c.toValue(v)
		if ttl != 0 && ttl <= time.Now().Unix() {
			return errors.New("data expired")
		}
		return nil
	})
	if err != nil {
		zap.S().Warn(err)
		return nil
	}
	return val
}

// toBytes
//  @Description:
//  @param val
//  @param extra 若为-1则仅用作key，无法转换为value
//  @return []byte
//
func (c *Bolt) toBytes(val interface{}, extra int64) []byte {
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
	case *models.AnimeEntity:
		buf.WriteByte(0x20)
		buf.Write(GobToBytes(value))
	case *models.AnimeSeason:
		buf.WriteByte(0x21)
		buf.Write(GobToBytes(value))
	case *models.AnimeEp:
		buf.WriteByte(0x22)
		buf.Write(GobToBytes(value))
	case *models.AnimeExtra:
		buf.WriteByte(0x23)
		buf.Write(GobToBytes(value))
	default:
		buf.WriteByte(0xFF)
		buf.Write(GobToBytes(value))
	}
	return buf.Bytes()
}

func (c *Bolt) toValue(data []byte) (val interface{}, extra int64) {
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
		val = &models.AnimeEntity{}
		GobToValue(data[9:], val)
	case 0x21:
		val = &models.AnimeSeason{}
		GobToValue(data[9:], val)
	case 0x22:
		val = &models.AnimeEp{}
		GobToValue(data[9:], val)
	case 0x23:
		val = &models.AnimeExtra{}
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

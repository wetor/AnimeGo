package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	"go.uber.org/zap"
	"time"
)

type Bolt struct {
	db *bolt.DB
}

func NewBolt() *Bolt {
	return &Bolt{}
}

func (c *Bolt) Open(path string) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		zap.S().Warn(err)
		return
	}

	c.db = db
}

func (c *Bolt) Close() {
	err := c.db.Close()
	if err != nil {
		zap.S().Warn(err)
		return
	}
}

func (c *Bolt) Add(bucket string) {
	err := c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return nil
	})
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

func (c *Bolt) Get(bucket string, key, val interface{}) error {
	var ttl int64
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get(c.toBytes(key, -1))
		if v == nil {
			return errors.New("不存在")
		}
		ttl = c.toValue(v, val)
		if ttl != 0 && ttl <= time.Now().Unix() {
			return errors.New("已过期")
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
	buf.Write(GobToBytes(val))
	return buf.Bytes()
}

func (c *Bolt) toValue(data []byte, val interface{}) (extra int64) {
	_ = data[8]
	extra = int64(binary.LittleEndian.Uint64(data[0:8]))
	GobToValue(data[8:], val)
	return extra
}

func GobToBytes(val interface{}) []byte {
	buf2 := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf2)
	err := enc.Encode(val)
	if err != nil {
		panic(err)
	}
	return buf2.Bytes()
}

func GobToValue(data []byte, val interface{}) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(val)
	if err != nil {
		panic(err)
	}
}

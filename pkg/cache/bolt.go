package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"github.com/wetor/AnimeGo/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"sync"

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
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Warn("打开bolt数据库失败")
		return
	}
	c.db = db
}

func (c *Bolt) Close() {
	err := c.db.Close()
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Warn("关闭bolt数据库失败")
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
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Warn("关闭bolt数据库失败")
		return
	}
}

func (c *Bolt) Put(bucket string, key, val interface{}, ttl int64) {
	if val == nil {
		return
	}
	var expire int64
	if ttl > 0 {
		expire = time.Now().Unix() + ttl
	} else {
		expire = 0
	}
	dbKey := c.toBytes(key, -1)
	dbVal := c.toBytes(val, expire)
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		return b.Put(dbKey, dbVal)
	})
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Warn("bolt添加数据失败")
		return
	}
}

func (c *Bolt) BatchPut(bucket string, key, val []interface{}, ttl int64) {
	if val == nil || len(key) != len(val) {
		return
	}
	var expire int64
	if ttl > 0 {
		expire = time.Now().Unix() + ttl
	} else {
		expire = 0
	}
	dbKeys := make([][]byte, len(key))
	dbVals := make([][]byte, len(key))
	wg := sync.WaitGroup{}
	wg.Add(len(key))
	for i := 0; i < len(key); i++ {
		go func(i int) {
			dbKeys[i] = c.toBytes(key[i], -1)
			dbVals[i] = c.toBytes(val[i], expire)
			wg.Done()
		}(i)
	}
	wg.Wait()
	err := c.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		for i := 0; i < len(key); i++ {
			err := b.Put(dbKeys[i], dbVals[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Warn("bolt添加数据失败")
		return
	}
}

func (c *Bolt) Get(bucket string, key, val interface{}) error {
	var ttl int64
	var dbVal []byte
	dbKey := c.toBytes(key, -1)
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		dbVal = b.Get(dbKey)
		if dbVal == nil {
			return errors.NewAniError("Key不存在")
		}
		return nil
	})
	if err != nil {
		return err
	}
	ttl = c.toValue(dbVal, val)
	if ttl != 0 && ttl <= time.Now().Unix() {
		return errors.NewAniError("Key已过期")
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
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Error("Gob Encode失败")
	}
	return buf2.Bytes()
}

func GobToValue(data []byte, val interface{}) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(val)
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Error("Gob Decode失败")
	}
}

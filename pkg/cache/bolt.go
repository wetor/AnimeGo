package cache

import (
	"bytes"
	"encoding/binary"
	jsoniter "github.com/json-iterator/go"
	"github.com/wetor/AnimeGo/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"sync"

	"go.uber.org/zap"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	_ = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		dbVal = b.Get(dbKey)
		return nil
	})
	if dbVal == nil {
		return errors.NewAniError("Key不存在")
	}
	ttl = c.toValue(dbVal, val)
	if ttl != 0 && ttl <= time.Now().Unix() {
		c.Delete(bucket, key)
		return errors.NewAniError("Key已过期")
	}
	return nil
}

func (c *Bolt) GetValue(bucket string, key string) (int64, string, error) {
	var dbVal []byte
	dbKey := []byte(key)
	_ = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		dbVal = b.Get(dbKey)
		return nil
	})
	if dbVal == nil {
		return 0, "", errors.NewAniError("Key不存在")
	}
	if len(dbVal) <= 8 {
		return 0, "", nil
	}
	ttl := int64(binary.LittleEndian.Uint64(dbVal[0:8]))
	if ttl != 0 && ttl <= time.Now().Unix() {
		c.Delete(bucket, key)
		return 0, "", errors.NewAniError("Key已过期")
	}
	val := string(dbVal[8:])
	return ttl, val, nil
}

// GetAll
//  @Description: 获取bucket所有kv数据
//  @receiver *Bolt
//  @param bucket string
//  @param tk interface{} key类型转换临时变量
//  @param tv interface{} value类型转换临时变量
//  @param fn func(k, v interface{})
//
func (c *Bolt) GetAll(bucket string, tk, tv interface{}, fn func(k, v interface{})) {
	_ = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))

		_ = b.ForEach(func(k, v []byte) error {
			_ = json.Unmarshal(k, tk)
			ttl := c.toValue(v, tv)
			if ttl != 0 && ttl <= time.Now().Unix() {
				return nil
			}
			fn(tk, tv)
			return nil
		})
		return nil
	})
}

func (c *Bolt) Delete(bucket string, key interface{}) {
	dbKey := c.toBytes(key, -1)
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Delete(dbKey)
		if err != nil {
			return err
		}
		return nil
	})
	errors.NewAniErrorD(err).TryPanic()
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
	data, err := json.Marshal(val)
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Error("Json Encode失败")
	}
	buf.Write(data)
	return buf.Bytes()
}

func (c *Bolt) toValue(data []byte, val interface{}) (extra int64) {
	_ = data[8]
	extra = int64(binary.LittleEndian.Uint64(data[0:8]))
	err := json.Unmarshal(data[8:], val)
	if err != nil {
		zap.S().Debug(errors.NewAniErrorD(err))
		zap.S().Error("Json Decode失败")
	}
	return extra
}

func (c *Bolt) ListBucket() []string {
	list := make([]string, 0, 8)
	_ = c.db.View(func(tx *bolt.Tx) error {
		err := tx.ForEach(func(nm []byte, b *bolt.Bucket) error {
			list = append(list, string(nm))
			return nil
		})
		return err
	})
	return list
}

func (c *Bolt) ListKey(bucket string) []string {
	list := make([]string, 0, 16)
	_ = c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.ForEach(func(k, v []byte) error {
			list = append(list, string(k))
			return nil
		})
		return err
	})
	return list
}

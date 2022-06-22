package cache

import (
	"GoBangumi/utils"
	"errors"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
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
		glog.Errorln(err)
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
		glog.Errorln(err)
		return
	}

}
func (c *Bolt) Close() {
	err := c.db.Close()
	if err != nil {
		glog.Errorln(err)
		return
	}
}
func (c *Bolt) Put(bucket Bucket, key, val interface{}, ttl int64) {
	var expire int64
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if ttl > 0 {
			expire = time.Now().Unix() + ttl
		} else {
			expire = 0
		}
		err := b.Put(utils.ToBytes(key, -1), utils.ToBytes(val, expire))
		return err
	})
	if err != nil {
		glog.Errorln(err)
		return
	}
}
func (c *Bolt) Get(bucket Bucket, key interface{}) interface{} {
	var val interface{}
	var ttl int64
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v := b.Get(utils.ToBytes(key, -1))
		val, ttl = utils.ToValue(v)
		if ttl != 0 && ttl <= time.Now().Unix() {
			return errors.New("data expired")
		}
		return nil
	})
	if err != nil {
		glog.Errorln(err)
		return nil
	}
	return val
}

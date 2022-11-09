// Package cache
// @Description: 缓存包，用来调用缓存组件
package cache

type Cache interface {
	Open(path string)
	Close()
	Add(bucket string)
	Put(bucket string, key, val interface{}, ttl int64)
	BatchPut(bucket string, key, val []interface{}, ttl int64)
	Get(bucket string, key, val interface{}) error
}

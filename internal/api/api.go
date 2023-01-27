package api

import (
	"context"

	"github.com/wetor/AnimeGo/internal/models"
)

type FilterManager interface {
	Update(ctx context.Context, items []*models.FeedItem)
}

type Cacher interface {
	CacheOpener
	CacheSetter
	CacheGetter
}

type DownloaderManagerCacheDeleter interface {
	DeleteCache(fullname string)
}

type CacheOpener interface {
	Open(file string)
	Close()
}

type CacheSetter interface {
	Add(bucket string)
	Put(bucket string, key, val interface{}, ttl int64)
	BatchPut(bucket string, key, val []interface{}, ttl int64)
	Delete(bucket string, key interface{}) error
}

type CacheGetter interface {
	Get(bucket string, key, val interface{}) error
	GetValue(bucket string, key interface{}) (int64, string, error)
	GetAll(bucket string, tk, tv interface{}, fn func(k, v interface{}))
	ListBucket() []string
	ListKey(bucket string) []string
}

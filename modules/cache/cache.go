package cache

type Bucket string

const (
	// DefaultBucket 默认bucket，储存任何信息
	DefaultBucket Bucket = "default_bucket"

	// RssBucket 订阅信息
	RssBucket Bucket = "rss_bucket"

	// BgmBucket 以bangumi id为key，主要存储番剧的各种信息
	BgmBucket Bucket = "bgm_bucket"

	// BgmEpBucket 以bangumi id为key，主要存储番剧ep信息
	BgmEpBucket Bucket = "bgm_ep_bucket"

	// TmdbBucket 以tmdb id为key，主要存储番剧季度信息
	TmdbBucket Bucket = "tmdb_bucket"
)

var buckets = []Bucket{DefaultBucket, RssBucket, BgmBucket, BgmEpBucket, TmdbBucket}

type Cache interface {
	Open(dir string)
	Close()
	Put(bucket Bucket, key, val interface{}, ttl int64)
	Get(bucket Bucket, key interface{}) interface{}
}

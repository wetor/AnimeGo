package cache

const (
	// DefaultBucket 默认bucket，储存任何信息
	DefaultBucket string = "default_bucket"
)

var buckets = []string{"rss_mikan", "mikan_bangumi", "name_tmdb", "tmdb_season", "bgm_info", "bgm_ep", "client_bangumi"}

type Cache interface {
	Open(dir string)
	Close()
	Put(bucket string, key, val interface{}, ttl int64)
	Get(bucket string, key interface{}) interface{}
}

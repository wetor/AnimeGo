package cache

const (
	// DefaultBucket 默认bucket，储存任何信息
	DefaultBucket       = "default_bucket"
	RssMikanBucket      = "rss_mikan"
	MikanBangumiBucket  = "mikan_bangumi"
	NameTmdbBucket      = "name_tmdb"
	TmdbSeasonBucket    = "tmdb_season"
	BgmInfoBucket       = "bgm_info"
	BgmEpBucket         = "bgm_ep"
	ClientBangumiBucket = "client_bangumi"
	ClientStateBucket   = "client_state"
)

var buckets = []string{
	RssMikanBucket,
	MikanBangumiBucket,
	NameTmdbBucket,
	TmdbSeasonBucket,
	BgmInfoBucket,
	BgmEpBucket,
	ClientBangumiBucket,
	ClientStateBucket,
}

type Cache interface {
	Open(dir string)
	Close()
	Put(bucket string, key, val interface{}, ttl int64)
	Get(bucket string, key interface{}) interface{}
}

package cache

import . "GoBangumi/internal/models"

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

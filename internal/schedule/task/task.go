package task

import (
	"sync"
	"time"

	"github.com/wetor/AnimeGo/internal/api"
)

var (
	DBDir            string
	BangumiCache     api.CacheOpener
	BangumiCacheLock *sync.Mutex
)

type Options struct {
	DBDir            string
	BangumiCache     api.CacheOpener
	BangumiCacheLock *sync.Mutex
}

type TaskInfo struct {
	Name  string
	RunAt time.Time
	Cron  string
}

type Task interface {
	Cron() string
	NextTime() time.Time
	Name() string
	Run(force bool)
}

func Init(opts *Options) {
	DBDir = opts.DBDir
	BangumiCache = opts.BangumiCache
	BangumiCacheLock = opts.BangumiCacheLock
}

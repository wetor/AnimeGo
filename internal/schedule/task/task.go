package task

import (
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/wetor/AnimeGo/internal/api"
)

var SecondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

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

func Init(opts *Options) {
	DBDir = opts.DBDir
	BangumiCache = opts.BangumiCache
	BangumiCacheLock = opts.BangumiCacheLock
}

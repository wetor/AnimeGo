package anidata

import (
	"sync"

	"github.com/wetor/AnimeGo/internal/api"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
)

var (
	Cache            mem.Memorizer
	CacheTime        map[string]int64
	BangumiCache     api.CacheGetter
	BangumiCacheLock *sync.Mutex

	RedirectMikan      string
	RedirectBangumi    string
	RedirectThemoviedb string
)

type Options struct {
	Cache            mem.Memorizer
	CacheTime        map[string]int64
	BangumiCache     api.CacheGetter
	BangumiCacheLock *sync.Mutex

	RedirectMikan      string
	RedirectBangumi    string
	RedirectThemoviedb string
}

func Init(opts *Options) {
	Cache = opts.Cache
	CacheTime = opts.CacheTime
	BangumiCache = opts.BangumiCache
	BangumiCacheLock = opts.BangumiCacheLock

	RedirectMikan = opts.RedirectMikan
	RedirectBangumi = opts.RedirectBangumi
	RedirectThemoviedb = opts.RedirectThemoviedb
}

package bangumi

import (
	"sync"

	"github.com/wetor/AnimeGo/internal/api"
	mem "github.com/wetor/AnimeGo/pkg/memorizer"
)

type Options struct {
	Cache            mem.Memorizer
	CacheTime        int64
	BangumiCache     api.CacheGetter
	BangumiCacheLock *sync.Mutex

	Host string
}

type Entity struct {
	ID      int    `json:"id"`      // Bangumi ID
	NameCN  string `json:"name_cn"` // 中文名
	Name    string `json:"name"`    // 原名
	Eps     int    `json:"eps"`     // 集数
	AirDate string `json:"airdate"` // 可空

	Type     int `json:"type"`
	Platform int `json:"platform"`
}

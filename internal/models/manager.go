package models

import "github.com/wetor/AnimeGo/internal/pkg/client"

type DownloadStatus struct {
	Hash     string              `json:"hash"`
	State    client.TorrentState `json:"state"`
	Path     []string            `json:"path"`      // 文件存储相对路径
	ExpireAt int64               `json:"expire_at"` // 过期时间

	Init       bool `json:"init"`       // 是否初始化
	Renamed    bool `json:"renamed"`    // 是否已重命名/移动
	Downloaded bool `json:"downloaded"` // 是否已下载完成
	Seeded     bool `json:"seeded"`     // 是否做种
	Scraped    bool `json:"scraped"`    // 是否已经完成搜刮
}

func (d DownloadStatus) Expire(now int64) bool {
	return d.ExpireAt > 0 && d.ExpireAt-now <= 0
}

func (d DownloadStatus) PathNull() bool {
	if d.Path == nil || len(d.Path) == 0 {
		return true
	}
	for _, r := range d.Path {
		if len(r) == 0 {
			return true
		}
	}
	return false
}

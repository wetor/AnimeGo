package models

type FeedItem struct {
	Url      string `json:"url"`      // 必要。Link，详情页连接，用于下一步解析番剧信息
	Name     string `json:"name"`     // 必要。种子名
	Date     string `json:"date"`     // 可选。发布日期
	Type     string `json:"type"`     // 可选。下载类型，[application/x-bittorrent]
	Download string `json:"download"` // 必要。下载链接
	Length   int64  `json:"length"`   // 可选。种子大小
}

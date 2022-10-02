package models

type FeedItem struct {
	Url     string `json:"url"`     // Link，详情页连接，用于下一步解析番剧信息
	Name    string `json:"name"`    // 种子名
	Date    string `json:"date"`    // 发布日期
	Torrent string `json:"torrent"` // 种子连接
	Hash    string `json:"hash"`    // 种子hash，唯一ID
}

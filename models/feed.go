package models

type FeedItem struct {
	Url     string // Link，详情页连接，用户下一步解析番剧信息
	Name    string // 种子名
	Date    string // 发布日期
	Torrent string // 种子连接
	Hash    string // 种子hash，唯一ID
}

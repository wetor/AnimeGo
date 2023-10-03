package models

type FeedItem struct {
	BangumiID     int    `json:"bangumi_id"`  // BangumiID和MikanUrl二选一。优先BangumiID
	MikanUrl      string `json:"mikan_url"`   // BangumiID和MikanUrl二选一。Mikan详情页连接，用于下一步解析番剧信息
	TorrentUrl    string `json:"torrent_url"` // 必要。下载链接
	Name          string `json:"name"`        // 可选。种子名
	Date          string `json:"date"`        // 可选。发布日期
	Type          string `json:"type"`        // 可选。下载类型，[application/x-bittorrent]
	Length        int64  `json:"length"`      // 可选。种子大小
	ParseOverride *AnimeParseOverride
}

package models

type FeedItem struct {
	Url      string `json:"url"`      // feed解析。Link，详情页连接，用于下一步解析番剧信息
	Name     string `json:"name"`     // feed解析。种子名
	Date     string `json:"date"`     // feed解析。发布日期
	Type     string `json:"type"`     // feed解析。下载类型，[application/x-bittorrent]
	Download string `json:"download"` // feed解析。下载链接
	Length   int64  `json:"length"`   // feed解析。种子大小

	Hash         string       `json:"hash"`          // bt唯一hash
	DownloadType string       `json:"download_type"` // 下载链接类型，[magnet, torrent]
	NameParsed   *TitleParsed `json:"parsed"`        // 标题解析信息
}

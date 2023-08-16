package models

type BaseDBEntity struct {
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
}

type AnimeDBEntity struct {
	BaseDBEntity `json:"info"`
	Init         bool `json:"init"`       // 是否初始化
	Renamed      bool `json:"renamed"`    // 是否已重命名/移动
	Downloaded   bool `json:"downloaded"` // 是否已下载完成
	Seeded       bool `json:"seeded"`     // 是否做种
	Scraped      bool `json:"scraped"`    // 是否已经完成搜刮
}

type EpisodeDBEntity struct {
	File string `json:"file"`
	Type int8   `json:"type"`
	Ep   int    `json:"ep"`
}

type SeasonDBEntity struct {
	BaseDBEntity `json:"info"`
	Season       int               `json:"season"`
	Episodes     []EpisodeDBEntity `json:"episodes"`
}

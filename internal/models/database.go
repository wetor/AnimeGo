package models

type BaseDBEntity struct {
	Hash     string `json:"hash"`
	Name     string `json:"name"`
	CreateAt int64  `json:"create_at"`
	UpdateAt int64  `json:"update_at"`
}

type AnimeDBEntity struct {
	BaseDBEntity `json:"info"`
}

type SeasonDBEntity struct {
	BaseDBEntity `json:"info"`
	Season       int `json:"season"`
}

type StateDB struct {
	Renamed bool `json:"renamed"` // 是否已重命名/移动
	Scraped bool `json:"scraped"` // 是否已经完成搜刮
}

type EpisodeDBEntity struct {
	BaseDBEntity
	StateDB
	Season int  `json:"season"`
	Type   int8 `json:"type"`
	Ep     int  `json:"ep"`
}

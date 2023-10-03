package models

import "fmt"

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
	Seeded     bool `json:"seeded"`     // 是否做种
	Downloaded bool `json:"downloaded"` // 是否已下载完成
	Renamed    bool `json:"renamed"`    // 是否已重命名/移动
	Scraped    bool `json:"scraped"`    // 是否已经完成搜刮
}

type EpisodeDBEntity struct {
	BaseDBEntity `json:"info"`
	StateDB      `json:"state"`
	Season       int         `json:"season"`
	Type         AnimeEpType `json:"type"`
	Ep           int         `json:"ep"`
}

func (e EpisodeDBEntity) Key() string {
	return fmt.Sprintf("E%d-%v", e.Ep, e.Type)
}

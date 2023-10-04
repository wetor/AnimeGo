package res

import (
	"github.com/wetor/AnimeGo/third_party/bangumi/model"
)

type Req struct {
	Keyword string `json:"keyword"`
	// 排序规则
	// match meilisearch 的默认排序，按照匹配程度
	// heat 收藏人数
	// rank 排名由高到低
	// score 评分
	Sort   string    `json:"sort"`
	Filter ReqFilter `json:"filter"`
}

type ReqFilter struct { //nolint:musttag
	Type    []model.SubjectType `json:"type"`     // or
	Tag     []string            `json:"tag"`      // and
	AirDate []string            `json:"air_date"` // and
	Score   []string            `json:"rating"`   // and
	Rank    []string            `json:"rank"`     // and
	// NSFW    null.Bool           `json:"nsfw"`
}

type ReponseSubject struct {
	Date string `json:"date"`
	// Image   string           `json:"image"`
	Type uint8 `json:"type"`
	// Summary string           `json:"summary"`
	Name   string `json:"name"`
	NameCN string `json:"name_cn"`
	// Tags    []res.SubjectTag `json:"tags"`
	// Score   float64          `json:"score"`
	ID model.SubjectIDType `json:"id"`
	// Rank uint32 `json:"rank"`
}

type SearchPaged struct {
	Data   []*ReponseSubject `json:"data"`
	Total  int64             `json:"total"`
	Limit  int               `json:"limit"`
	Offset int               `json:"offset"`
}

package bangumi

type Entity struct {
	ID      int    `json:"id"`      // Bangumi ID
	NameCN  string `json:"name_cn"` // 中文名
	Name    string `json:"name"`    // 原名
	Eps     int    `json:"eps"`     // 集数
	AirDate string `json:"airdate"` // 可空

	Type     int `json:"type"`
	Platform int `json:"platform"`
}

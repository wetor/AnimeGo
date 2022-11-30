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

type Ep struct {
	Ep      int    `json:"ep"`      // 当前集，从下载文件名解析
	AirDate string `json:"airdate"` // 当前集播放日期，从bgm获取
	ID      int    `json:"id"`      // 当前集bgm id
	//Duration string // 当前集时长
	//Desc     string // 当前集简介
	//Name     string // 当前集标题
	//NameCN   string // 当前集中文标题

	SubjectID int     `json:"subject_id"`
	Sort      float32 `json:"sort"`
	Type      int     `json:"type"`
}

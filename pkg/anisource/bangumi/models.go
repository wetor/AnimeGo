package bangumi

type Entity struct {
	ID      int    // Bangumi ID
	NameCN  string // 中文名
	Name    string // 原名
	Eps     int    // 集数
	AirDate string // 可空
}

type Ep struct {
	Ep       int    // 当前集，从下载文件名解析
	Date     string // 当前集播放日期，从bgm获取
	Duration string // 当前集时长
	EpDesc   string // 当前集简介
	EpName   string // 当前集标题
	EpNameCN string // 当前集中文标题
	EpID     int    // 当前集bgm id
}

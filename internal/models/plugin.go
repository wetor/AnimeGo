package models

type TitleParsed struct {
	TitleRaw   string `json:"title_raw"`
	Name       string `json:"title_jp"`
	NameCN     string `json:"title_zh"`
	NameEN     string `json:"title_en"`
	Season     int    `json:"season"`
	SeasonRaw  string `json:"season_raw"`
	Ep         int    `json:"episode"`
	Sub        string `json:"sub"`
	Group      string `json:"group"`
	Definition string `json:"resolution"`
	Source     string `json:"source"`
}

type Plugin struct {
	Enable bool           `json:"enable"`
	Type   string         `json:"type"`
	File   string         `json:"file"`
	Code   *string        `json:"code"`
	Args   map[string]any `json:"args"`
	Vars   map[string]any `json:"vars"`
}

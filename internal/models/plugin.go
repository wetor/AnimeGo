package models

type Plugin struct {
	Enable bool           `json:"enable"`
	Type   string         `json:"type"`
	File   string         `json:"file"`
	Args   map[string]any `json:"args"`
	Vars   map[string]any `json:"vars"`
}

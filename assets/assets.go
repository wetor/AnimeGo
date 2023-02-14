package assets

import (
	"embed"
)

var (
	//go:embed plugin
	//go:embed plugin/filter/Auto_Bangumi/__init__.py
	Plugin embed.FS
)

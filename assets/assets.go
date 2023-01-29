package assets

import (
	"embed"
)

var (
	//go:embed plugin
	//go:embed plugin/filter/pylib/__init__.py
	Plugin embed.FS
)

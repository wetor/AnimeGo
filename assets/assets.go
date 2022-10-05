package assets

import (
	"embed"
)

var (
	//go:embed plugin
	Plugin embed.FS
	//go:embed config
	Config embed.FS
)

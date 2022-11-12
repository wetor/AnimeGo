package dmhy

import "github.com/wetor/AnimeGo/internal/models"

type DmhyAdapter struct {
	ThemoviedbKey string
}

func (adapter DmhyAdapter) Parse(opt *models.AnimeParseOptions) *models.AnimeEntity {
	return ParseDmhy(opt.Name, opt.Date, adapter.ThemoviedbKey)
}

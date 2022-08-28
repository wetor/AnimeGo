package mikan

import "AnimeGo/internal/models"

type MikanAdapter struct {
	ThemoviedbKey string
}

func (adapter MikanAdapter) Parse(opt *models.AnimeParseOptions) *models.AnimeEntity {
	return ParseMikan(opt.Name, opt.Url, adapter.ThemoviedbKey)
}

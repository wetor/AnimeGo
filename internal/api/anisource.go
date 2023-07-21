package api

import "github.com/wetor/AnimeGo/internal/models"

type AniSource interface {
	Parse(opt *models.AnimeParseOptions) *models.AnimeEntity
}

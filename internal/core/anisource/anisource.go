package anisource

import (
	"GoBangumi/internal/models"
)

type AniSource interface {
	Parse(opt *models.AnimeParseOptions) *models.AnimeEntity
}

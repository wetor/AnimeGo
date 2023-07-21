package api

import "github.com/wetor/AnimeGo/internal/models"

type ParserPlugin interface {
	Parse(string) *models.TitleParsed
}

type ParserManager interface {
	Parse(*models.ParseOptions) *models.AnimeEntity
}

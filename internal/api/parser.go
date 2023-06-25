package api

import "github.com/wetor/AnimeGo/internal/models"

type ParserPlugin interface {
	Parse(string) (*models.TitleParsed, error)
}

type ParserManager interface {
	Parse(*models.ParseOptions) (*models.AnimeEntity, error)
}

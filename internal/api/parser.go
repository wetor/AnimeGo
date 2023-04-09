package api

import "github.com/wetor/AnimeGo/internal/models"

type ParserPlugin interface {
	Parse(string) *models.TitleParsed
}

type ParserManager interface {
	Parse(title, url string) *models.AnimeEntity
}

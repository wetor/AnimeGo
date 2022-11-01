package parser

import (
	"github.com/wetor/AnimeGo/internal/models"
	"github.com/wetor/AnimeGo/pkg/errors"
	"github.com/wetor/AnimeGo/third_party/poketto"
)

func ParseTitle(title string) (*models.ParseResult, error) {
	parse := poketto.NewEpisode(title)
	parse.TryParse()
	if parse.ParseErr != nil {
		return nil, errors.NewAniErrorD(parse.ParseErr)
	}
	return &models.ParseResult{
		Ep:         parse.Ep,
		Definition: parse.Definition,
		Subtitle:   parse.Sub,
		Source:     parse.Source,
	}, nil
}

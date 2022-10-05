package parser

import (
	"AnimeGo/internal/models"
	"AnimeGo/third_party/poketto"
)

func ParseTitle(title string) (*models.ParseResult, error) {
	parse := poketto.NewEpisode(title)
	parse.TryParse()
	if parse.ParseErr != nil {
		return nil, parse.ParseErr
	}
	return &models.ParseResult{
		Ep:         parse.Ep,
		Definition: parse.Definition,
		Subtitle:   parse.Sub,
		Source:     parse.Source,
	}, nil
}
